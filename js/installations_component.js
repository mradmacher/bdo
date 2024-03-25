import { openModal, closeModal } from "./modal_helpers.js"

function formatCode(code, dangerous) {
  let formattedCode = [
    code.slice(0, 2),
    code.slice(2, 4),
    code.slice(4, 6),
  ].join(" ")
  if (dangerous) {
    formattedCode += "*"
  }

  return formattedCode;
}

class CodeTemplate {
  constructor() {
    this.template = document.getElementById('code-template');
  }

  build(code) {
    let element = this.template.content.cloneNode(true);
    element.querySelector('.code-slot').textContent = code;

    return element;
  }
}

class InstallationCapabilityTemplate {
  constructor() {
    this.template = document.getElementById('installation-capability-template');
  }

  build(capability) {
    let element = this.template.content.cloneNode(true);
    let formattedCode = [
      capability.WasteCode.slice(0, 2),
      capability.WasteCode.slice(2, 4),
      capability.WasteCode.slice(4, 6),
    ].join(" ")
    if (capability.Dangerous) {
      formattedCode += "*";
    }
    element.querySelector('.waste-code').textContent = formattedCode;
    element.querySelector('.process-code').textContent = capability.ProcessCode;
    element.querySelector('.quantity').textContent = capability.Quantity;

    let processStatus = 'teal';
    if (capability.ProcessCode.startsWith('D')) {
      processStatus = 'orange';
    }
    element.querySelector('.process-code').classList.add(processStatus);

    let quantityStatus = 'olive';
    if (capability.Quantity > 1500) {
      quantityStatus = 'purple';
    }
    element.querySelector('.quantity').classList.add(quantityStatus);

    return element;
  }
}

class InstallationTemplate {
  constructor() {
    this.template = document.getElementById('installation-template');
  }

  build(installation, onShowDetails) {
    let element = this.template.content.cloneNode(true);
    element.querySelector('.name-slot').textContent = installation.Name;
    element.querySelector('.address-slot').textContent = `${installation.Address.Line1},  ${installation.Address.Line2}`;
    element.querySelector('.show-details-action').addEventListener('click', (event) => {
      onShowDetails(installation);
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
      element.querySelector('.waste-codes-slot').append(new CodeTemplate().build(code));
    })
    processCodes.sort().forEach((code, i) => {
      element.querySelector('.process-codes-slot').append(new CodeTemplate().build(code));
    })

    return element;
  }
}

export class InstallationDetailsView {
  constructor() {
    this.modal = document.querySelector('.modal.installation-details');
    this.modal.querySelectorAll('.button.cancel').forEach((elem) => {
      elem.addEventListener('click', (event) => {
        closeModal(this.modal);
      })
    })
  }

  hide() {
    closeModal(this.modal);
  }

  show(installation) {
    this.modal.querySelector('.installation-name').textContent = installation.Name
    this.modal.querySelector('.installation-address').textContent = `${installation.Address.Line1},  ${installation.Address.Line2}`
    this.modal.querySelector('.installation-capabilities').innerHTML = '';
    installation.Capabilities.forEach((capability, i) => {
      this.modal.querySelector('.installation-capabilities').append(
        new InstallationCapabilityTemplate().build(capability)
      )
    })
    openModal(this.modal);
  }
}

export class InstallationsComponent {
  constructor(elementId) {
    this.element = document.getElementById(elementId);
    this.detailsView = new InstallationDetailsView();
  }

  clear() {
    this.element.innerHTML = '';
  }

  showDetails(installation) {
    this.detailsView.show(installation);
  }

  addInstallation(installation) {
    this.element.appendChild(
      new InstallationTemplate().build(installation, (installation) => {
        this.showDetails(installation);
      })
    )
  }
}
