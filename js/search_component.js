import { WasteHinter, ProcessHinter } from "./hinters.js"

export class SearchComponent {
  constructor(elementId, onReset, onSearch) {
    this.element = document.getElementById(elementId)
    this.onReset = onReset
    this.onSearch = onSearch
    this.element.querySelector('form').reset()
    this.setWasteHint(null, null, null)
    this.setProcessHint(null)
    this.processHinter = new ProcessHinter();
    this.wasteHinter = new WasteHinter();
    this.searched = false;

    this.element.querySelector('form [type="reset"]').addEventListener('click', (event) => {
      this.searched = false;
      this.setWasteHint(null, null, null)
      this.setProcessHint(null)
      this.onReset()
    })

    this.element.querySelector('form [type="submit"]').addEventListener('click', (event) => {
      event.preventDefault()

      this.searched = true;
      this.onSearch(this.collectSearchParams())
      this.setHints()
    })
  }

  repeatSearch() {
    if (this.searched) {
      this.onSearch(this.collectSearchParams())
    }
  }

  collectSearchParams() {
    let params = {}
    let paramValue = this.element.querySelector('[name=waste]').value
    if(paramValue) {
      params["wc"] = paramValue
    }
    paramValue = this.element.querySelector('[name=process]').value
    if(paramValue) {
      params["pc"] = paramValue
    }
    paramValue = this.element.querySelector('[name=state]').value
    if(paramValue) {
      params["sc"] = paramValue
    }

    return params
  }

  setHints() {
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
      descA = this.wasteHinter.descriptionFor(codeA);
    }
    if(codeB) {
      descB = this.wasteHinter.descriptionFor(codeB);
    }
    if(codeC) {
      descC = this.wasteHinter.descriptionFor(codeC);
    }
    this.setWasteHint(descA, descB, descC)

    let processInputElement = this.element.querySelector('[name=process]')
    let process = processInputElement.value
    if(process) {
      process = process.toUpperCase()
      processInputElement.value = process
      this.setProcessHint(this.processHinter.descriptionFor(process))
    } else {
      this.setProcessHint(null)
    }
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
