export class InstallationsComponent {
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
