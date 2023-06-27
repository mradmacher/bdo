class InstallationsView {
  constructor(selector) {
    this.selector = selector
    this.element = $(selector)
  }

  clear() {
    this.element.html('')
  }

  addInstallation(installation) {
    let template = $($('#installation-template').html())
    template.find('.name').text(installation.name)
    template.find('.address').text(installation.address)
    installation.capabilities.forEach((capability, i) => {
      let capabilityTemplate = $($('#installation-capability-template').html())
      capabilityTemplate.find('.waste-code').text(capability.wasteCode)
      capabilityTemplate.find('.process-code').text(capability.processCode)
      capabilityTemplate.find('.quantity').text(capability.quantity)

      let processStatus = 'teal'
      if (capability.processCode.startsWith('D')) {
        processStatus = 'orange'
      }
      capabilityTemplate.find('.process-code').addClass(processStatus)

      let quantityStatus = 'olive'
      if (capability.quantity > 1500) {
        quantityStatus = 'purple'
      }
      capabilityTemplate.find('.quantity').addClass(quantityStatus)

      template.find('.capabilities').append(capabilityTemplate)
    })
    this.element.append(template)
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
      installations.forEach((installation, i) => {
        installationsView.addInstallation(installation)
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
