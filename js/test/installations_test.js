import test from 'ava'
import { JSDOM } from 'jsdom'
import {InstallationsView} from '../app.js'

test.before(t => {
  t.context.html =
    `
      <template id="installation-capability-template">
        <div class="capability"</div>
          <span class="waste-code"></span>
          <span class="process-code"></span>
          <span class="quantity"></span>
        </div>
      </template>

      <template id="installation-template">
        <div class="installation">
          <div class="name"></div>
          <div class="address"></div>
          <div class="capabilities">
          </div>
        </div>
      </template>
      <div id="installations"></div>
    `
})

test('adds installation with capabilities', t => {
  const dom = new JSDOM(t.context.html)
  global.document = dom.window.document
  let view = new InstallationsView('installations')
  view.addInstallation({
    Name: "Test",
    Address: {
      Line1: "Address Line 1",
      Line2: "Address Line 2",
    },
    Capabilities: [
      {
        WasteCode: "010101",
        ProcessCode: "R12",
        Quantity: "1000",
      },
      {
        WasteCode: "020202",
        Dangerous: true,
        ProcessCode: "D10",
        Quantity: "500",
      }
    ],

  })
  let installationsElement = document.querySelectorAll(".installation")
  t.deepEqual(installationsElement.length, 1)

  let installationElement = installationsElement[0]
  t.deepEqual(installationElement.querySelector(".name").textContent, "Test")
  t.regex(installationElement.querySelector(".address").textContent, /Address Line 1/)
  t.regex(installationElement.querySelector(".address").textContent, /Address Line 2/)

  let capabilitiesElement = installationElement.querySelectorAll(".capability")
  t.deepEqual(capabilitiesElement.length, 2)

  let capabilityElement = capabilitiesElement[0]
  t.deepEqual(capabilityElement.querySelector(".waste-code").textContent, "01 01 01")
  t.deepEqual(capabilityElement.querySelector(".process-code").textContent, "R12")
  t.deepEqual(capabilityElement.querySelector(".quantity").textContent, "1000")

  capabilityElement = capabilitiesElement[1]
  t.deepEqual(capabilityElement.querySelector(".waste-code").textContent, "02 02 02*")
  t.deepEqual(capabilityElement.querySelector(".process-code").textContent, "D10")
  t.deepEqual(capabilityElement.querySelector(".quantity").textContent, "500")
})
