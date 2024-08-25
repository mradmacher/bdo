function getHrefContent(href) {
  return new Promise((resolve, reject) => {
    if (href.startsWith("#")) {
      let elementId = href.substring(1);
      resolve(document.getElementById(elementId).content.cloneNode(true));
    } else {
      axios.get(href).then((response) => {
        resolve(response.data);
      }).catch((error) => {
        reject(error.response.data)
      })
    }
  })
}

export function initModalTriggers(element) {
  element.querySelectorAll("[data-modal-target]").forEach((actionElement) => {
    actionElement.addEventListener('click', (event) => {
      event.preventDefault();
      let modalFragment = document.querySelector(actionElement.getAttribute("data-modal-source")).content.cloneNode(true);
      let wrapper = document.querySelector(actionElement.getAttribute("data-modal-target"));
      modalFragment.querySelector("[data-modal-title]").textContent = actionElement.getAttribute("data-modal-title");
      let href = actionElement.getAttribute("href");
      let modalBody = modalFragment.querySelector("[data-modal-body]")
      getHrefContent(href).then((content) => {
        if (content instanceof DocumentFragment) {
          modalFragment.querySelector("[data-modal-body]").append(content);
        } else {
          modalFragment.querySelector("[data-modal-body]").innerHTML = content;
        }
        wrapper.append(modalFragment);
        let modal = wrapper.querySelector("[data-modal]");
        wrapper.dispatchEvent(new Event("contentLoaded"));
        openModal(modal);
      })
    })
  })
}

export function openModal(modal) {
  modal.querySelectorAll("[data-close-modal]").forEach((closeElement) => {
    closeElement.addEventListener('click', (event) => {
      modal.classList.remove('is-active');
      modal.remove();
    })
  })
  modal.classList.add('is-active');
}

export function closeModal(modal) {
  modal.classList.remove('is-active');
}
