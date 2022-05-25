import {ref, unref} from "vue";
import cloneDeep from 'lodash-es/cloneDeep';

export function useForm(modelRef, rulesRef) {
  console.log('useForm', modelRef.value)

  const initialModel = cloneDeep(unref(modelRef));
  const validateInfos = ref({})

  const validate = () => {
    let success = true

    const model = unref(modelRef)
    const rules = unref(rulesRef)
    const ruleKeys = unref(Object.keys(rules))

    ruleKeys.forEach((key, index) => {
      const errorMap = {}
      let pass = true
      if (!errorMap[key]) errorMap[key] = []

      rules[key].forEach((item, index) => {
        if (item.required) pass &&= checkRequired(key, item, model, errorMap)
        if (pass && item.email) pass &&= checkEmail(key, item, model, errorMap)
        if (pass && item.regex) pass &&= checkRegex(key, item, model, errorMap)

        success &&= pass
      })

      validateInfos.value[key] = errorMap
    })

    return success;
  }

  const reset = () => {
    console.log('reset')
    modelRef.value = {...initialModel}
  }

  return {
    validate, reset, validateInfos,
  };
}

export function checkRequired(key, item, model, errMap) {
  let pass = true;
  if (item.required && !model[key]) pass = false

  if (!pass) errMap[key].push(item.msg)
  return pass;
}

export function checkEmail(key, item, model, errMap) {
  const regx=/^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/;
  const pass = regx.test(model[key]);

  if (!pass) errMap[key].push(item.msg)
  return pass;
}

export function checkRegex(key, item, model, errMap) :any {
  const regx = new RegExp(item.regex);
  const pass = regx.test(model[key]);

  if (!pass) errMap[key].push(item.msg)
  return pass;
}

