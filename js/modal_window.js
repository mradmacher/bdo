export class ModalWindow extends HTMLElement {
  constructor() {
    super();
    this.classList.add("modal");
    let templateContent = document.getElementById("modal-window-template").content;
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.appendChild(templateContent.cloneNode(true));

    shadowRoot.querySelectorAll("[data-close-button]").forEach((closeElement) => {
      closeElement.addEventListener("click", (event) => {
        this.close();
      })
    })
  }

  open() {
    this.classList.add("is-active");
  }

  close() {
    this.classList.remove("is-active");
  }

  setBody(content, callback) {
    let bodyElement = this.shadowRoot.querySelector("slot[name=\"body\"]");
    bodyElement.innerHTML = content;
    callback(bodyElement);
  }

  set innerHTML(content) {
    this.shadowRoot.querySelector("slot[name=\"body\"]").innerHTML = content;
  }
}
