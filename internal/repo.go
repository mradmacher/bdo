package bdo

import (
	"errors"
	"fmt"
	"strings"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "bdo"

type RecordNotFound struct {
	Collection string
	Id         int64
}

func (e RecordNotFound) Error() string {
	return fmt.Sprintf("Could not find record with %d ID in %s collection", e.Id, e.Collection)
}

func IsRecordNotFound(err error) bool {
	_, ok := err.(RecordNotFound)
	return ok
}

type Repository struct {
	DBUri string
	db *sql.DB
}

type DbRecord interface {
	Add(r *Repository) (int64, error)
}

type Material struct {
	Id   int64
	Code string
}

type Capability struct {
	Id           int64
	WasteCode    string
	Dangerous    bool
	ProcessCode  string
	ActivityCode string
	Quantity     int
	Materials    []Material
}

type Address struct {
	Line1     string
	Line2     string
	StateCode string
	Lat       string
	Lng       string
}

type Installation struct {
	Id           int64
	Name         string
	Address      Address
	Capabilities []Capability
}

type SearchParams map[string]string

func (r *Repository) Connect() error {
	var err error
	uri := r.DBUri
	if uri == "" {
		return errors.New("Set database URI")
	}
	r.db, err = sql.Open("sqlite3", uri)
	if err != nil {
		return errors.Join(errors.New("Problem opening the database"), err)
	}

	return nil
}

func (r *Repository) Disconnect() error {
	return r.db.Close()
}

func (r *Repository) Purge() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM materials")
	if err != nil {
		return errors.Join(errors.New("Problem removing materials"), err)
	}
	_, err = tx.Exec("DELETE FROM capabilities")
	if err != nil {
		return errors.Join(errors.New("Problem removing capabilities"), err)
	}
	_, err = tx.Exec("DELETE FROM installations")
	if err != nil {
		return errors.Join(errors.New("Problem removing installations"), err)
	}
	tx.Commit()
	return nil
}

func (inst Installation) Add(r *Repository) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	result, err := tx.Exec(
		"INSERT INTO installations (name, address_line1, address_line2, state_code, lat, lng) values (?, ?, ?, ?, ?, ?)",
		inst.Name, inst.Address.Line1, inst.Address.Line2, inst.Address.StateCode, inst.Address.Lat, inst.Address.Lng,
	)
	if err != nil {
		return 0, errors.Join(errors.New("Executing INSERT INTO installations failed"), err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Join(errors.New("Could not obtain installation last inserted ID"), err)
	}
	cstmt, err := tx.Prepare(
		"INSERT INTO capabilities (installation_id, waste_code, dangerous, process_code, activity_code, quantity) values (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return 0, errors.Join(errors.New("Preparing capabilities statement failed"), err)
	}
	defer cstmt.Close()

	mstmt, err := tx.Prepare(
		"INSERT INTO materials (capability_id, code) values (?, ?)",
	)
	if err != nil {
		return 0, errors.Join(errors.New("Preparing materials statement failed"), err)
	}
	defer mstmt.Close()

	for _, capability := range inst.Capabilities {
		result, err = cstmt.Exec(id, capability.WasteCode, capability.Dangerous, capability.ProcessCode, capability.ActivityCode, capability.Quantity)
		if err != nil {
			return 0, err
		}
		cid, err := result.LastInsertId()
		if err != nil {
			return 0, errors.Join(errors.New("Could not obtain capability last inserted ID"), err)
		}
		for _, material := range capability.Materials {
			_, err = mstmt.Exec(cid, material.Code)
			if err != nil {
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) LoadInstallationCapabilities(id int64) ([]Capability, error) {
	var capabilities []Capability

	query := "SELECT id, waste_code, dangerous, process_code, activity_code, quantity FROM capabilities WHERE installation_id = ?"

	rows, err := r.db.Query(query, id)
	defer rows.Close()

	for rows.Next() {
		var c Capability
		err := rows.Scan(&c.Id, &c.WasteCode, &c.Dangerous, &c.ProcessCode, &c.ActivityCode, &c.Quantity)
		if err != nil {
			return nil, errors.Join(errors.New("Reading capabilities failed"), err)
		}
		capabilities = append(capabilities, c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Join(errors.New("Not all capabilities read"), err)

	}

	return capabilities, nil
}

func (r *Repository) SearchCapabilities(id int64, params SearchParams) ([]Capability, error) {
	var capabilities []Capability
	var whereCond []string
	var whereArgs []any

	whereCond = append(whereCond, "installation_id = ?")
	whereArgs = append(whereArgs, id)

	for k, v := range params {
		switch k {
		case "waste_code":
			whereCond = append(whereCond, "waste_code = ?")
			whereArgs = append(whereArgs, v)
		}
	}

	query := "SELECT id, waste_code, dangerous, process_code, activity_code, quantity FROM capabilities"
	var whereClause string
	if len(whereCond) > 0 {
		whereClause = strings.Join(whereCond, " AND ")
		query = query + " WHERE " + whereClause
	}

	rows, err := r.db.Query(query, whereArgs...)
	defer rows.Close()

	for rows.Next() {
		var c Capability
		err := rows.Scan(&c.Id, &c.WasteCode, &c.Dangerous, &c.ProcessCode, &c.ActivityCode, &c.Quantity)
		if err != nil {
			return nil, errors.Join(errors.New("Reading capabilities failed"), err)
		}
		mquery := "SELECT id, code FROM materials WHERE capability_id = ?"
		mrows, err := r.db.Query(mquery, c.Id)
		defer mrows.Close()
		for mrows.Next() {
			var m Material
			err := mrows.Scan(&m.Id, &m.Code)
			if err != nil {
				return nil, errors.Join(errors.New("Reading materials failed"), err)
			}
			c.Materials = append(c.Materials, m)
		}
		capabilities = append(capabilities, c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Join(errors.New("Not all capabilities read"), err)

	}

	return capabilities, nil
}

func (r *Repository) Summarize(params SearchParams) ([]Capability, error) {
	var capabilities []Capability
	var whereCond []string
	var whereArgs []any

	for k, v := range params {
		switch k {
		case "process_code":
			whereCond = append(whereCond, "process_code = ?")
			whereArgs = append(whereArgs, v)
		case "waste_code":
			whereCond = append(whereCond, "waste_code = ?")
			whereArgs = append(whereArgs, v)
		case "state_code":
			whereCond = append(whereCond, "installation_id IN (SELECT id FROM installations WHERE state_code = ?)")
			whereArgs = append(whereArgs, v)
		}
	}

	query := "SELECT waste_code, dangerous, process_code, sum(quantity) FROM capabilities"
	var whereClause string
	if len(whereCond) > 0 {
		whereClause = strings.Join(whereCond, " AND ")
		query = query + " WHERE " + whereClause
	}
	query = query + " GROUP BY waste_code, dangerous, process_code"

	rows, err := r.db.Query(query, whereArgs...)
	if err != nil {
		return nil, errors.Join(errors.New("Searching capabilities failed"), err)
	}
	defer rows.Close()
	for rows.Next() {
		c := Capability{}
		err = rows.Scan(&c.WasteCode, &c.Dangerous, &c.ProcessCode, &c.Quantity)
		if err != nil {
			return nil, errors.Join(errors.New("Scaning capability failed"), err)
		}
		capabilities = append(capabilities, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Join(errors.New("Not all capabilities read"), err)
	}

	return capabilities, nil
}

func (r *Repository) Search(params SearchParams) ([]*Installation, error) {
	var installations []*Installation
	var where_cond []string
	var where_args []any

	for k, v := range params {
		switch k {
		case "process_code":
			where_cond = append(where_cond, "process_code = ?")
			where_args = append(where_args, v)
		case "waste_code":
			where_cond = append(where_cond, "waste_code = ?")
			where_args = append(where_args, v)
		case "state_code":
			where_cond = append(where_cond, "state_code = ?")
			where_args = append(where_args, v)
		}
	}
	query := "SELECT id, name, address_line1, address_line2, state_code, lat, lng from installations"
	var where_clause string
	if len(where_cond) > 0 {
		where_clause = strings.Join(where_cond, " AND ")
		query = query + " WHERE id IN (SELECT DISTINCT installation_id FROM capabilities WHERE " + where_clause + ")"
	}

	rows, err := r.db.Query(query, where_args...)
	if err != nil {
		return nil, errors.Join(errors.New("Searching installations failed"), err)
	}
	defer rows.Close()
	for rows.Next() {
		inst := Installation{}
		err = rows.Scan(&inst.Id, &inst.Name, &inst.Address.Line1, &inst.Address.Line2, &inst.Address.StateCode, &inst.Address.Lat, &inst.Address.Lng)
		if err != nil {
			return nil, errors.Join(errors.New("Scaning installation failed"), err)
		}
		installations = append(installations, &inst)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Join(errors.New("Not all installations read"), err)
	}

	for _, inst := range installations {
		inst.Capabilities, err = r.LoadInstallationCapabilities(inst.Id)
		if err != nil {
			return nil, errors.Join(errors.New("Loading capabilities failed"), err)
		}
	}

	return installations, nil
}

func (r *Repository) Find(id int64, inst *Installation) error {
	row := r.db.QueryRow("SELECT id, name, address_line1, address_line2, state_code, lat, lng FROM installations WHERE id = ?", id)
	err := row.Scan(&inst.Id, &inst.Name, &inst.Address.Line1, &inst.Address.Line2, &inst.Address.StateCode, &inst.Address.Lat, &inst.Address.Lng)
	if err != nil {
		if err == sql.ErrNoRows {
			return RecordNotFound{Collection: "installations", Id: id}
		} else {
			return errors.Join(errors.New("Finding installation failed"), err)
		}
	}
	inst.Capabilities, err = r.LoadInstallationCapabilities(inst.Id)
	if err != nil {
		return errors.Join(errors.New("Loading capabilities failed"), err)
	}
	return nil
}
