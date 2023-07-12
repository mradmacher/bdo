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

export class SearchView {
  constructor(selector, installationsView, mapView) {
    this.selector = selector
    this.element = $(selector)
    this.installationsView = installationsView
    this.mapView = mapView
    this.element.find('.waste-hint.code-a').hide()
    this.element.find('.waste-hint.code-b').hide()
    this.element.find('.waste-hint.code-c').hide()
    this.element.find('.process-hint').hide()
    this.element.trigger('form reset')

    this.element.find('.button[type="reset"]').click((event) => {
      installationsView.clear()
      mapView.clear()
      this.element.find('.waste-hint.code-a').hide()
      this.element.find('.waste-hint.code-b').hide()
      this.element.find('.waste-hint.code-c').hide()
      this.element.find('.process-hint').hide()
      this.element.trigger('form reset')
    })

    this.element.find('.button[type="submit"]').click((event) => {
      event.preventDefault()
      installationsView.clear()
      mapView.clear()

      let params = {
        'wc': this.element.find('[name=waste]').val(),
        'pc': this.element.find('[name=process]').val(),
        'sc': this.element.find('[name=state]').val(),
      }

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

      $.ajax({
        method: "GET",
        url: "/api/installations",
        data: params,
        dataType: "json",
      }).done(function(installations) {
        installations.forEach((installation, i) => {
          installationsView.addInstallation(installation)
          mapView.addInstallation(installation)
        })
      })
      let descA
      let descB
      let descC
      let code = $('.ui.form .field [name=waste]').val()
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

      if(descA) {
        this.element.find('.waste-hint.code-a').text(descA)
        this.element.find('.waste-hint.code-a').show()
      } else {
        this.element.find('.waste-hint.code-a').text("Nieznany kod")
      }
      if(descB) {
        this.element.find('.waste-hint.code-b').text(descB)
        this.element.find('.waste-hint.code-b').show()
      } else {
        this.element.find('.waste-hint.code-b').hide()
      }
      if(descC) {
        this.element.find('.waste-hint.code-c').text(descC)
        this.element.find('.waste-hint.code-c').show()
      } else {
        this.element.find('.waste-hint.code-c').hide()
      }

      let value
      let process = this.element.find('.field [name=process]').val()
      if(process) {
        process = process.toUpperCase()
        this.element.find('form .field [name=process]').val(process)
        value = processDescs[process]
      }
      if(value) {
        this.element.find('.process-hint').text(value)
        this.element.find('.process-hint').show()
      } else {
        this.element.find('.process-hint').text("Nieznany kod")
      }
    })
  }

  setProcess(code, description) {
    this.element.find('form input[name="process"]').val(code)
    this.element.find('.process-hint').text(description)
    this.element.find('.process-hint').show()
  }

  setWaste(code, descA, descB, descC) {
    this.element.find('form input[name="waste"]').val(code)
    this.element.find('.waste-hint.code-a').show()
    this.element.find('.waste-hint.code-a').text(descA)
    this.element.find('.waste-hint.code-b').show()
    this.element.find('.waste-hint.code-b').text(descB)
    this.element.find('.waste-hint.code-c').show()
    this.element.find('.waste-hint.code-c').text(descC)
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
