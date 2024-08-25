import test from 'ava'
import { JSDOM } from 'jsdom'
import {SearchComponent} from '../search_component.js'

test.before(t => {
  t.context.html =
    `
      <div id="search">
        <form>
          <input type="text" name="waste">
          <div class="waste-hint code-a"></div>
          <div class="waste-hint code-b"></div>
          <div class="waste-hint code-c"></div>
          <input type="text" name="process">
          <div class="process-hint"></div>
          <select class="ui search selection simple dropdown" name="state">
            <option value=""></option>
            <option value="02">A</option>
            <option value="04">B</option>
            <option value="06">C</option>
          </select>
          <button type="submit"></button>
          <button type="reset"></button>
        </form>
      </div>
    `
})

test('returns search params on search', t => {
  const dom = new JSDOM(t.context.html)
  global.document = dom.window.document
  let onReset = () => {}
  let searchedParams
  let onSearch = (params) => {
    searchedParams = params
  }
  let search = new SearchComponent('search', onReset, onSearch)

  document.querySelector("[name=waste]").value = "101010"
  document.querySelector("[name=process]").value = "R1"
  document.querySelector("[name=state]").value = "06"
  document.querySelector('form [type="submit"]').click()

  t.assert(Object.keys(searchedParams).includes('wc'))
  t.deepEqual(searchedParams["wc"], "101010")
  t.assert(Object.keys(searchedParams).includes('pc'))
  t.deepEqual(searchedParams["pc"], "R1")
  t.assert(Object.keys(searchedParams).includes('sc'))
  t.deepEqual(searchedParams["sc"], "06")
})

test('returns nothing when search params not entered', t => {
  const dom = new JSDOM(t.context.html)
  global.document = dom.window.document
  let onReset = () => {}
  let searchedParams
  let onSearch = (params) => {
    searchedParams = params
  }
  let search = new SearchComponent('search', onReset, onSearch)

  document.querySelector('form [type="submit"]').click()

  t.deepEqual(searchedParams, {})
})
