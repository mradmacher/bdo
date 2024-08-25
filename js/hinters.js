import { wasteCodes, wasteCodeDescs } from "./waste_catalog.js"
import { processCodes, processCodeDescs } from "./process_catalog.js"

export class WasteHinter {
  relatedCodesFor(code) {
    if (code == '') {
      return wasteCodes['00'];
    } else {
      return wasteCodes[code];
    }
  }

  descriptionFor(code) {
    return wasteCodeDescs[code.replace("*", "")];
  }
}

export class ProcessHinter {
  relatedCodesFor(code) {
    if (code == '') {
      return processCodes;
    } else {
      return null;
    }
  }

  descriptionFor(code) {
    return processCodeDescs[code];
  }
}
