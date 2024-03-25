import test from 'ava';
import { JSDOM } from 'jsdom';
import { InstallationsComponent } from '../installations_component.js';

test.before(t => {
  t.context.html =
    `
      <template id="code-template">
        <span class="code-slot"></span>
      </template>

      <template id="installation-capability-template">
        <div class="capability"</div>
          <span class="waste-code"></span>
          <span class="process-code"></span>
          <span class="quantity"></span>
        </div>
      </template>

      <template id="installation-template">
        <div class="installation">
          <div class="name-slot"></div>
          <div class="address-slot"></div>
          <div class="waste-codes-slot"></div>
          <div class="process-codes-slot"></div>
          <button class="show-details-action"></button>
        </div>
      </template>

      <div id="installations"></div>
      <div class="modal installation-details"></div>
    `
})

test('adds installation with capabilities', t => {
  const dom = new JSDOM(t.context.html)
  global.document = dom.window.document
  let view = new InstallationsComponent('installations')
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
  t.deepEqual(installationElement.querySelector(".name-slot").textContent, "Test")
  t.regex(installationElement.querySelector(".address-slot").textContent, /Address Line 1/)
  t.regex(installationElement.querySelector(".address-slot").textContent, /Address Line 2/)

  let wasteCodesElement = installationElement.querySelectorAll(".waste-codes-slot > .code-slot")
  t.deepEqual(wasteCodesElement.length, 2)

  let processCodesElement = installationElement.querySelectorAll(".process-codes-slot > .code-slot")
  t.deepEqual(processCodesElement.length, 2)

  t.deepEqual(wasteCodesElement[0].textContent, "01 01 01")
  t.deepEqual(wasteCodesElement[1].textContent, "02 02 02*")

  t.deepEqual(processCodesElement[0].textContent, "D10")
  t.deepEqual(processCodesElement[1].textContent, "R12")
})
