
export function getElectron(): boolean {
    const userAgent = navigator.userAgent.toLowerCase()
    console.log(`userAgent ${userAgent}`)

    if (userAgent.indexOf('electron') > -1){
        return true
    }

    return false
}

export function isWindows(): boolean {
    const agent = navigator.userAgent.toLowerCase();
    // console.log('agent', agent)
    if (agent.indexOf("win") > -1 || agent.indexOf("wow") > -1) {
       return true
    }

    return false
}

export function proxyArrToVal(proxyArr: any[]): any[] {
    const items = []as any[]
    proxyArr.forEach((item) => {
        items.push(item)
    })

    return items
}