import { codes, codeDescs } from "./waste_catalog.js"
import { processes, processDescs } from "./process_catalog.js"
import {MapComponent} from "./map_component.js"
import {SearchComponent} from "./search_component.js"
import {InstallationsComponent} from "./installations_component.js"


function openModal($el) {
  $el.classList.add('is-active');
}

function closeModal($el) {
  $el.classList.remove('is-active');
}

export function updateUrlSearchParams(params) {
  if ('URLSearchParams' in window) {
    let searchParams = new URLSearchParams();
    let searchParamsProvided = false;
    for(let p in params) {
      if(params[p]) {
        searchParams.set(p, params[p]);
        searchParamsProvided = true;
      }
    }
    let newRelativePathQuery = window.location.pathname;
    if(searchParamsProvided) {
      newRelativePathQuery = newRelativePathQuery + '?' + searchParams.toString();
    }
    history.pushState(null, '', newRelativePathQuery);
  }
}

export class InstallationRequest {
  constructor() {
    this.url = "/api/installations"
  }

  search(params) {
    return new Promise((resolve, reject) => {
      axios.get(this.url, {
        params: params
      }).then(function(response) {
        console.log(response)
        resolve(response.data)
      }).catch(function(error) {
        reject(error.response.data)
      })
    })
  }
/*
  search(params) {
    axios.get(this.url, {
      params: params
    })
    .then(function(installations) {
      resolve(installations)
    })
    .catch(function(error) {
      console.log(error)
      reject(error);
    })
  }
  search(params) {
    return new Promise((resolve, reject) => {
      $.ajax({
        method: "GET",
        url: this.url,
        data: params,
        dataType: "json",
      }).done(function(installations) {
        resolve(installations)
      }).fail(function(xhr, status, error) {
        reject(xhr.responseJSON.errors);
      })
    })
  }
*/
}

class CodeDescRowBuilder {
  static build(code, desc, callback) {
    var template = document.getElementById('code-desc-row-template').content.cloneNode(true);
    template.querySelector('.code').textContent = code;
    template.querySelector('.description').textContent = desc;
    template.querySelector('.action').addEventListener('click', (event) => {
      callback(code, desc)
    })
    return template
  }
}

class CodeDescHeaderRowBuilder {
  static build(code, desc) {
    let template = document.getElementById('code-desc-header-row-template').content.cloneNode(true);
    template.querySelector('.code').textContent = code;
    template.querySelector('.description').textContent = desc;

    return template
  }
}

export class ProcessSelectorView {
  constructor() {
    this.modal = document.querySelector('.modal.processes');
    this.modal.querySelector('.process-list').content = '';
    this.modal.querySelectorAll('.button.cancel').forEach((elem) => {
      elem.addEventListener('click', (event) => {
        closeModal(this.modal);
      })
    })
  }

  hide() {
    closeModal(this.modal);
  }

  show() {
    openModal(this.modal);
  }

  load(onSelect) {
    for (let code in processDescs) {
      let template = CodeDescRowBuilder.build(code, processDescs[code], (code, desc) => {
        onSelect(code, desc)
      })
      this.modal.querySelector('.process-list').append(template);
    }
  }
}

export class WasteSelectorView {
  constructor() {
    this.modal = document.querySelector('.modal.wastes');
    this.modal.querySelectorAll('.selected-waste-header').forEach((elem) => {
      elem.remove();
    })
    this.modal.querySelectorAll('.button.cancel').forEach((elem) => {
      elem.addEventListener('click', (event) => {
        closeModal(this.modal);
      })
    })
    this.modal.querySelector('.waste-list').innerHTML = '';
    this.selectedDescs = []
    this.selectedCode = '00'
  }

  select(code, desc) {
    this.modal.querySelector('.waste-list-header').append(
      CodeDescHeaderRowBuilder.build(code, desc)
    )
  }

  hide() {
    closeModal(this.modal);
  }

  show() {
    openModal(this.modal)
  }

  load(onSelect) {
    this.modal.querySelector('.waste-list').innerHTML = '';
    codes[this.selectedCode].forEach((code, i) => {
      let template = CodeDescRowBuilder.build(code, codeDescs[code.replace("*", "")], (code, desc) => {
        this.selectedCode = code
        this.selectedDescs.push(desc)
        onSelect(code, desc)
      })
      this.modal.querySelector('.waste-list').append(template)
    })
  }
}

document.addEventListener("DOMContentLoaded", () => {
  let googleMapsApiKey = document.getElementById('google-maps-api-key').getAttribute('data-value');
  let installationsComponent
  let searchComponent
  let mapComponent = new MapComponent(googleMapsApiKey)

  mapComponent.initMap("map").then(()=> {
    installationsComponent = new InstallationsComponent('installations')
    searchComponent = new SearchComponent('search',
      () => {
        installationsComponent.clear()
        mapComponent.clear()
      },
      (params) => {
        updateUrlSearchParams(params)

        installationsComponent.clear()
        mapComponent.clear()

        new InstallationRequest().search(params)
          .then((installations) => {
            installations.forEach((installation, i) => {
              installationsComponent.addInstallation(installation)
              mapComponent.addInstallation(installation)
            })
          })
      }
    )

    document.querySelector('.search.process').addEventListener('click', (event) => {
      event.preventDefault();
      let processSelectorView = new ProcessSelectorView
      processSelectorView.load((code, desc) => {
        searchComponent.setProcess(code, desc)
        processSelectorView.hide()
      })
      processSelectorView.show()
    })

    document.querySelector('.search.waste').addEventListener('click', (event) => {
      event.preventDefault();
      var wasteListView = new WasteSelectorView
      wasteListView.load((code, desc) => {
        wasteListView.select(code, desc)
        wasteListView.load((code, desc) => {
          wasteListView.select(code, desc)
          wasteListView.load((code, desc) => {
            searchComponent.setWaste(wasteListView.selectedCode, ...wasteListView.selectedDescs)
            wasteListView.hide()
          })
        })
      })
      wasteListView.show()
    })
  })

})
