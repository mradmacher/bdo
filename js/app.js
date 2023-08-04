import { codes, codeDescs } from "./waste_catalog.js"
import { processes, processDescs } from "./process_catalog.js"

export function updateUrlSearchParams(params) {
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
}

export class InstallationRequest {
  constructor() {
    this.url = "/api/installations"
  }

  search(params) {
    return new Promise((resolve, reject) => {
      axios.get(this.url, {
        params: params
      }).then(function(response) {
        console.log(response)
        resolve(response.data)
      }).catch(function(error) {
        reject(error.response.data)
      })
    })
  }
/*
  search(params) {
    axios.get(this.url, {
      params: params
    })
    .then(function(installations) {
      resolve(installations)
    })
    .catch(function(error) {
      console.log(error)
      reject(error);
    })
  }
  search(params) {
    return new Promise((resolve, reject) => {
      $.ajax({
        method: "GET",
        url: this.url,
        data: params,
        dataType: "json",
      }).done(function(installations) {
        resolve(installations)
      }).fail(function(xhr, status, error) {
        reject(xhr.responseJSON.errors);
      })
    })
  }
*/
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
