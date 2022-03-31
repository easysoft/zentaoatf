export function isWindows(): boolean {
    const agent = navigator.userAgent.toLowerCase();
    // console.log('agent', agent)
    if (agent.indexOf("win") > -1 || agent.indexOf("wow") > -1) {
       return true
    }

    return false
}