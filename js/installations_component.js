function formatCode(code, dangerous) {
  let formattedCode = [
    code.slice(0, 2),
    code.slice(2, 4),
    code.slice(4, 6),
  ].join(" ")
  if (dangerous) {
    formattedCode += "*"
  } else {
    formattedCode += " "
  }
  return formattedCode;
}
export class InstallationsComponent {
  constructor(elementId) {
    this.element = document.getElementById(elementId);
  }

  clear() {
    this.element.innerHTML = '';
  }

  showDetails(installation) {
    let modal = document.querySelector('.modal.installation-details');
    modal.querySelector('.installation-name').textContent = installation.Name
    modal.querySelector('.installation-address').textContent = `${installation.Address.Line1},  ${installation.Address.Line2}`
    modal.querySelector('.installation-capabilities').innerHTML = '';
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

      modal.querySelector('.installation-capabilities').append(capabilityTemplate)
    })
    modal.querySelectorAll('.button.cancel').forEach((elem) => {
      elem.addEventListener('click', (event) => {
        modal.classList.remove('is-active');
      })
    })
    modal.classList.add('is-active');
  }

  addInstallation(installation) {
    let template = document.getElementById('installation-template').content.cloneNode(true)
    template.querySelector('.name').textContent = installation.Name
    template.querySelector('.address').textContent = `${installation.Address.Line1},  ${installation.Address.Line2}`
    template.querySelector('.action.show-details').addEventListener('click', (event) => {
      this.showDetails(installation);
    })
    let wasteCodes = []
    let processCodes = []
    installation.Capabilities.forEach((capability, i) => {
      let formattedCode = formatCode(capability.WasteCode, capability.Dangerous);
      if (!wasteCodes.includes(formattedCode)) {
        wasteCodes.push(formattedCode);
      }
      if (!processCodes.includes(capability.ProcessCode)) {
        processCodes.push(capability.ProcessCode);
      }
    });
    wasteCodes.sort().forEach((code, i) => {
      let wasteCodeTemplate = document.getElementById('code-template').content.cloneNode(true);
      wasteCodeTemplate.querySelector('.code').textContent = code;
      template.querySelector('.waste-codes').append(wasteCodeTemplate);
    })
    processCodes.sort().forEach((code, i) => {
      let processCodeTemplate = document.getElementById('code-template').content.cloneNode(true)
      processCodeTemplate.querySelector('.code').textContent = code;
      template.querySelector('.process-codes').append(processCodeTemplate)
    })

    this.element.append(template)
  }
}
