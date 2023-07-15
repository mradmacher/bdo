export class MapView {
  constructor(apiKey) {
    (g=>{var h,a,k,p="The Google Maps JavaScript API",c="google",l="importLibrary",q="__ib__",m=document,b=window;b=b[c]||(b[c]={});var d=b.maps||(b.maps={}),r=new Set,e=new URLSearchParams,u=()=>h||(h=new Promise(async(f,n)=>{await (a=m.createElement("script"));e.set("libraries",[...r]+"");for(k in g)e.set(k.replace(/[A-Z]/g,t=>"_"+t[0].toLowerCase()),g[k]);e.set("callback",c+".maps."+q);a.src=`https://maps.${c}apis.com/maps/api/js?`+e;d[q]=f;a.onerror=()=>h=n(Error(p+" could not load."));a.nonce=m.querySelector("script[nonce]")?.nonce||"";m.head.append(a)}));d[l]?console.warn(p+" only loads once. Ignoring:",g):d[l]=(f,...n)=>r.add(f)&&u().then(()=>d[l](f,...n))})({
      key: apiKey,
      v: "weekly",
      // Use the 'v' parameter to indicate the version to use (weekly, beta, alpha, etc.).
      // Add other bootstrap parameters as needed, using camel case.
    });

    this.markers = []
    this.mapCenter = { lat: 52.24, lng: 21.00 }
    this.mapZoom = 6
  }

  clear() {
    this.map.setZoom(this.mapZoom)
    this.map.setCenter(this.mapCenter)
    while (this.markers.length) {
      this.markers.pop().setMap(null)
    }
  }

  addInstallation(installation) {
    const latLng = new google.maps.LatLng(
      installation.Address.Lat,
      installation.Address.Lng
    )
    let marker = new google.maps.Marker({
      map: this.map,
      position: latLng,
      label: installation.Name,
      title: `${installation.Name}\n${installation.Address.Line1}\n${installation.Address.Line2}`,
    })
    this.markers.push(marker)
  }

  async initMap(elementId) {
    const { Map } = await google.maps.importLibrary("maps");

    this.map = new Map(document.getElementById(elementId), {
      center: this.mapCenter,
      zoom: this.mapZoom,
    });
  }
}

export class InstallationsView {
  constructor(elementId) {
    this.element = document.getElementById(elementId)
  }

  clear() {
    this.element.innerHTML = ''
  }

  addInstallation(installation) {
    let template = document.getElementById('installation-template').content.cloneNode(true)
    template.querySelector('.name').textContent = installation.Name
    template.querySelector('.address').textContent = `${installation.Address.Line1},  ${installation.Address.Line2}`
    installation.Capabilities.forEach((capability, i) => {
      let capabilityTemplate = document.getElementById('installation-capability-template').content.cloneNode(true)
      let formattedCode = [
        capability.WasteCode.slice(0, 2),
        capability.WasteCode.slice(2, 4),
        capability.WasteCode.slice(4, 6),
      ].join(" ")
      if (capability.Dangerous) {
        formattedCode += "*"
      }
      capabilityTemplate.querySelector('.waste-code').textContent = formattedCode
      capabilityTemplate.querySelector('.waste-code').textContent = formattedCode
      capabilityTemplate.querySelector('.process-code').textContent = capability.ProcessCode
      capabilityTemplate.querySelector('.quantity').textContent = capability.Quantity

      let processStatus = 'teal'
      if (capability.ProcessCode.startsWith('D')) {
        processStatus = 'orange'
      }
      capabilityTemplate.querySelector('.process-code').classList.add(processStatus)

      let quantityStatus = 'olive'
      if (capability.Quantity > 1500) {
        quantityStatus = 'purple'
      }
      capabilityTemplate.querySelector('.quantity').classList.add(quantityStatus)

      template.querySelector('.capabilities').append(capabilityTemplate)
    })
    this.element.append(template)
  }
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
}

export class SearchView {
  constructor(elementId, onReset, onSearch) {
    this.element = document.getElementById(elementId)
    this.onReset = onReset
    this.onSearch = onSearch
    this.element.querySelector('form').reset()
    this.setWasteHint(null, null, null)
    this.setProcessHint(null)

    this.element.querySelector('form [type="reset"]').addEventListener('click', (event) => {
      this.setWasteHint(null, null, null)
      this.setProcessHint(null)
      this.onReset()
    })

    this.element.querySelector('form [type="submit"]').addEventListener('click', (event) => {
      event.preventDefault()

      let params = {
        'wc': this.element.querySelector('[name=waste]').value,
        'pc': this.element.querySelector('[name=process]').value,
        'sc': this.element.querySelector('[name=state]').value,
      }

      this.onSearch(params)

      let descA
      let descB
      let descC
      let code = this.element.querySelector('[name=waste]').value
      let codeA = code.slice(0, 2)
      let codeB = code.slice(0, 4)
      let codeC = code.slice(0, 6)
      if(codeA.length < 2) {
        codeA = ''
      }
      if(codeB.length < 4) {
        codeB = ''
      }
      if(codeC.length < 6) {
        codeC = ''
      }
      if(codeA) {
        descA = codeDescs[codeA]
      }
      if(codeB) {
        descB = codeDescs[codeB]
      }
      if(codeC) {
        descC = codeDescs[codeC]
      }
      this.setWasteHint(descA, descB, descC)

      let processInputElement = this.element.querySelector('.field [name=process]')
      let process = processInputElement.value
      if(process) {
        process = process.toUpperCase()
        processInputElement.value = process
        this.setProcessHint(processDescs[process])
      } else {
        this.setProcessHint(null)
      }
    })
  }

  setWasteHint(hintA, hintB, hintC) {
    let hintElement = this.element.querySelector('.waste-hint.code-a')
    if(hintA) {
      hintElement.textContent = hintA
      hintElement.hidden = false
      hintElement.style.display = "block"
    } else {
      hintElement.textContent = ''
      hintElement.hidden = true
      hintElement.style.display = "none"
    }
    hintElement = this.element.querySelector('.waste-hint.code-b')
    if(hintB) {
      hintElement.textContent = hintB
      hintElement.hidden = false
      hintElement.style.display = "block"
    } else {
      hintElement.textContent = ''
      hintElement.hidden = true
      hintElement.style.display = "none"
    }
    hintElement = this.element.querySelector('.waste-hint.code-c')
    if(hintC) {
      hintElement.textContent = hintC
      hintElement.hidden = false
      hintElement.style.display = "block"
    } else {
      hintElement.textContent = ''
      hintElement.hidden = true
      hintElement.style.display = "none"
    }
  }

  setProcessHint(hint) {
    let hintElement = this.element.querySelector('.process-hint')
    if(hint) {
      hintElement.textContent = hint
      hintElement.hidden = false
      hintElement.style.display = "block"
    } else {
      hintElement.textContent = ''
      hintElement.hidden = true
      hintElement.style.display = "none"
    }
  }

  setProcess(code, description) {
    this.element.querySelector('form input[name="process"]').value = code
    this.setProcessHint(description)
  }

  setWaste(code, descA, descB, descC) {
    this.element.querySelector('form input[name="waste"]').value = code
    this.setWasteHint(descA, descB, descC)
  }
}

class CodeDescRowBuilder {
  static build(code, desc, callback) {
    var template = $($('#code-desc-row-template').html())
    template.find('.code').text(code)
    template.find('.description').text(desc)
    template.click((event) => {
      callback(code, desc)
    })
    return template
  }
}

class CodeDescHeaderRowBuilder {
  static build(code, desc) {
    let template = $($('#code-desc-header-row-template').html())
    template.find('.code').text(code)
    template.find('.description').text(desc)

    return template
  }
}

export class ProcessSelectorView {
  constructor() {
    this.modal = $('.ui.modal.processes')
    this.modal.find('.process-list').html('')
  }

  hide() {
    this.modal.modal('hide')
  }

  show() {
    this.modal.modal('show')
  } load(onSelect) { for (let code in processDescs) { let template = CodeDescRowBuilder.build(code, processDescs[code], (code, desc) => {
        onSelect(code, desc)
      })
      this.modal.find('.process-list').append(template)
    }
  }
}

export class WasteListView {
  constructor() {
    this.modal = $('.ui.modal.wastes')
    this.modal.find('.selected-waste-header').remove()
    this.modal.find('.waste-list').html('')
    this.selectedDescs = []
    this.selectedCode = '00'
  }

  select(code, desc) {
    this.modal.find('.waste-list-header').append(
      CodeDescHeaderRowBuilder.build(code, desc)
    )
  }

  hide() {
    this.modal.modal('hide')
  }

  show() {
    this.modal.modal('show')
  }

  load(onSelect) {
    this.modal.find('.waste-list').html('')
    codes[this.selectedCode].forEach((code, i) => {
      let template = CodeDescRowBuilder.build(code, codeDescs[code.replace("*", "")], (code, desc) => {
        this.selectedCode = code
        this.selectedDescs.push(desc)
        onSelect(code, desc)
      })
      this.modal.find('.waste-list').append(template)
    })
  }
}
