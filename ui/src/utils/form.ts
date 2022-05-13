import {ref, unref} from "vue";

export function checkRequired(item) :any {
  if (!item.required) return {pass: true};

  return {key: 'required', pass: false, msg: item.msg};
}

export function useForm(modelRef, rulesRef) {
  const validateInfos = ref({})

  const validate = () => {
    let success = true;

    const rules = unref(rulesRef)
    const ruleKeys = unref(Object.keys(rules))

    ruleKeys.forEach((key, index) => {
      const errorMap = {}
      rules[key].forEach((item, index) => {
        const { key, pass, msg } = checkRequired(item)
        if (!pass) {
          success = false

          if (!errorMap[key]) errorMap[key] = []
          errorMap[key].push(msg)
        }
      })

      validateInfos.value[key] = errorMap
    })

    return success;
  }

  const reset = () => {
    console.log('reset')
  }

  return {
    validate, reset, validateInfos,
  };
}
