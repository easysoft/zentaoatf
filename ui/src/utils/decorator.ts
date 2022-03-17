import {disableStatusMap} from "@/utils/const";

export function disableStatusDef(disabled: boolean): string {
  const str = disabled ? '0' : '1'

  console.log('sdfdsf', disableStatusMap)

  let ret = 'N/A'
  disableStatusMap.forEach((item) => {
    if (item.value === str) {
      ret = item.label
    }
  })

  return ret
}