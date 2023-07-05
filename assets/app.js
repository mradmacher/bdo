class InstallationsView {
  constructor(selector, map) {
    this.selector = selector
    this.element = $(selector)
    this.map = map
    this.markers = []
  }

  clear() {
    this.element.html('')
    this.map.setZoom(6)
    this.map.setCenter(new google.maps.LatLng(52.24, 21.00))
    while (this.markers.length) {
      this.markers.pop().setMap(null)
    }
  }

  addInstallation(installation) {
    let template = $($('#installation-template').html())
    template.find('.name').text(installation.Name)
    template.find('.address').text(`${installation.Address.Line1},  ${installation.Address.Line2}`)
    installation.Capabilities.forEach((capability, i) => {
      let capabilityTemplate = $($('#installation-capability-template').html())
      let formattedCode = [
        capability.WasteCode.slice(0, 2),
        capability.WasteCode.slice(2, 4),
        capability.WasteCode.slice(4, 6),
      ].join(" ")
      if (capability.Dangerous) {
        formattedCode += "*"
      }
      capabilityTemplate.find('.waste-code').text(formattedCode)
      capabilityTemplate.find('.process-code').text(capability.ProcessCode)
      capabilityTemplate.find('.quantity').text(capability.Quantity)

      let processStatus = 'teal'
      if (capability.ProcessCode.startsWith('D')) {
        processStatus = 'orange'
      }
      capabilityTemplate.find('.process-code').addClass(processStatus)

      let quantityStatus = 'olive'
      if (capability.Quantity > 1500) {
        quantityStatus = 'purple'
      }
      capabilityTemplate.find('.quantity').addClass(quantityStatus)

      template.find('.capabilities').append(capabilityTemplate)
    })
    this.element.append(template)

    const latLng = new google.maps.LatLng(installation.Address.Lat, installation.Address.Lng)
    let marker = new google.maps.Marker({
      map: this.map,
      position: latLng,
      label: installation.Name,
      title: `${installation.Name}\n${installation.Address.Line1}\n${installation.Address.Line2}`,
    })
    this.markers.push(marker)
  }
}

class SearchView {
  constructor(selector, installationsView) {
    this.selector = selector
    this.element = $(selector)
    this.installationsView = installationsView
    this.element.find('.waste-hint.code-a').hide()
    this.element.find('.waste-hint.code-b').hide()
    this.element.find('.waste-hint.code-c').hide()
    this.element.find('.process-hint').hide()
    this.element.trigger('form reset')

    this.element.find('.button[type="reset"]').click((event) => {
      installationsView.clear()
      this.element.find('.waste-hint.code-a').hide()
      this.element.find('.waste-hint.code-b').hide()
      this.element.find('.waste-hint.code-c').hide()
      this.element.find('.process-hint').hide()
      this.element.trigger('form reset')
    })

    this.element.find('.button[type="submit"]').click((event) => {
      event.preventDefault()
      installationsView.clear()

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

class ProcessSelectorView {
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

class WasteListView {
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
