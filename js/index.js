import { WasteHinter } from "./waste_catalog.js"
import { ProcessHinter } from "./process_catalog.js"
import { MapComponent } from "./map_component.js"
import { SearchComponent } from "./search_component.js"
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
    this.indexUrl = "/api/installations"
    this.showUrl = "/api/installation"
  }

  search(params) {
    return new Promise((resolve, reject) => {
      axios.get(this.indexUrl, {
        params: params
      }).then(function(response) {
        resolve(response.data)
      }).catch(function(error) {
        reject(error.response.data)
      })
    })
  }

  show(id) {
    return new Promise((resolve, reject) => {
      axios.get(`${this.showUrl}/${id}`)
      .then(function(response) {
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

export class CodeSelectorView {
  constructor(selector, hinter, onSelect) {
    this.modal = document.querySelector(selector);
    this.hinter = hinter;
    this.modal.querySelectorAll('.button.cancel').forEach((elem) => {
      elem.addEventListener('click', (event) => {
        closeModal(this.modal);
      })
    })
    this.modal.querySelectorAll('.button.accept').forEach((elem) => {
      elem.addEventListener('click', (event) => {
        onSelect(this.selectedCode, this.selectedDescs);
        closeModal(this.modal);
      })
    })
    this.modal.querySelectorAll('.selected-header').forEach((elem) => {
      elem.remove();
    })
    this.modal.querySelector('.list').innerHTML = '';
    this.disableAccept();
    this.selectedDescs = [];
    this.selectedCode = '';
  }

  select(code, desc) {
    this.selectedCode = code
    this.selectedDescs.push(desc)
    this.modal.querySelector('.list-header').append(
      new CodeDescHeaderRowTemplate().build(code, desc)
    )
  }

  enableAccept() {
    this.modal.querySelector('.button.accept').removeAttribute('disabled');
  }

  disableAccept() {
    this.modal.querySelector('.button.accept').setAttribute('disabled', '');
  }

  hide() {
    closeModal(this.modal);
  }

  show() {
    this.load();
    openModal(this.modal);
  }

  load() {
    let relatedCodes = this.hinter.relatedCodesFor(this.selectedCode);

    this.modal.querySelector('.list').innerHTML = '';
    if (!relatedCodes) {
      this.enableAccept();
      return;
    }
    this.hinter.relatedCodesFor(this.selectedCode).forEach((code, i) => {
      let element = new CodeDescRowTemplate().build(code, this.hinter.descriptionFor(code), (code, desc) => {
        this.select(code, desc);
        this.load();
      })
      this.modal.querySelector('.list').append(element)
    })
  }
}

export class ProcessSelectorView extends CodeSelectorView {
  constructor(selector, onSelect) {
    super(selector, new ProcessHinter(), onSelect);
  }
}

export class WasteSelectorView extends CodeSelectorView {
  constructor(selector, onSelect) {
    super(selector, new WasteHinter(), onSelect);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  let googleMapsApiKey = document.getElementById('google-maps-api-key').getAttribute('data-value');
  let searchComponent
  let mapComponent = new MapComponent(googleMapsApiKey)

  mapComponent.initMap("map").then(()=> {
    searchComponent = new SearchComponent('search',
      () => {
        mapComponent.clear()
      },
      (params) => {
        updateUrlSearchParams(params)

        mapComponent.clear()

        new InstallationRequest().search(params)
          .then((installations) => {
            let listElement = document.getElementById('installations')
            listElement.innerHTML = installations;
            listElement.querySelectorAll('[data-installation]').forEach((installationElement) => {
              mapComponent.addInstallation({
                addressLat: installationElement.getAttribute('data-lat'),
                addressLng: installationElement.getAttribute('data-lng'),
                name: installationElement.querySelector('[data-name]').textContent,
                addressLine1: installationElement.querySelector('[data-address-line1]').textContent,
                addressLine2: installationElement.querySelector('[data-address-line2]').textContent,
              });

            })

            listElement.querySelectorAll('[data-show-details]').forEach((actionElement) => {
              actionElement.addEventListener('click', (event) => {
                let id = event.target.getAttribute('data-id');
                new InstallationRequest().show(id).then((result) => {
                  let detailsElement = document.querySelector('.installation-details');
                  detailsElement.innerHTML = result;
                  detailsElement.querySelectorAll('.button.cancel').forEach((cancelElement) => {
                    cancelElement.addEventListener('click', (event) => {
                      closeModal(detailsElement);
                    })
                  })
                  openModal(detailsElement);
                })
              })
            })
          })
      }
    )

    document.querySelector('.search.process').addEventListener('click', (event) => {
      event.preventDefault();
      new ProcessSelectorView('.modal.processes', (code, descs) => {
        searchComponent.setProcess(code, descs[0]);
      }).show();
    })

    document.querySelector('.search.waste').addEventListener('click', (event) => {
      event.preventDefault();
      new WasteSelectorView('.modal.wastes', (code, descs) => {
        searchComponent.setWaste(code, ...descs);
      }).show();
    })
  })

})
