export function initModalTriggers(element) {
  element.querySelectorAll("[data-modal-window]").forEach((actionElement) => {
    actionElement.addEventListener('click', (event) => {
      event.preventDefault();
      let modal = element.querySelector(`#${actionElement.getAttribute("data-modal-window")}`);
      let href = actionElement.getAttribute("data-href");
      if (href) {
        axios.get(href).then((response) => {
          modal.setBody(response.data, (bodyElement) => {
            initModalTriggers(bodyElement);
          })
          modal.open();
        })
      } else {
        modal.open();
      }
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
