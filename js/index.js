import { codes, codeDescs } from "./waste_catalog.js"
import { processes, processDescs } from "./process_catalog.js"
import { MapComponent } from "./map_component.js"
import { SearchComponent } from "./search_component.js"
import { InstallationsComponent } from "./installations_component.js"
import { openModal, closeModal } from "./modal_helpers.js"

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
}

class CodeDescRowTemplate {
  constructor() {
    this.template = document.getElementById('code-desc-row-template');
  }

  build(code, desc, onSelect) {
    var element = this.template.content.cloneNode(true);
    element.querySelector('.code-slot').textContent = code;
    element.querySelector('.description-slot').textContent = desc;
    element.querySelector('.select-action').addEventListener('click', (event) => {
      onSelect(code, desc);
    })

    return element;
  }
}

class CodeDescHeaderRowTemplate {
  constructor() {
    this.template = document.getElementById('code-desc-header-row-template');
  }

  build(code, desc) {
    let element = this.template.content.cloneNode(true);
    element.querySelector('.code').textContent = code;
    element.querySelector('.description').textContent = desc;

    return element;
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
      let element = new CodeDescRowTemplate().build(code, processDescs[code], (code, desc) => {
        onSelect(code, desc)
      })
      this.modal.querySelector('.process-list').append(element);
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
      new CodeDescHeaderRowTemplate().build(code, desc)
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
      let element = new CodeDescRowTemplate().build(code, codeDescs[code.replace("*", "")], (code, desc) => {
        this.selectedCode = code
        this.selectedDescs.push(desc)
        onSelect(code, desc)
      })
      this.modal.querySelector('.waste-list').append(element)
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
