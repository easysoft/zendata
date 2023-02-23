export function getElectron() {
    const userAgent = navigator.userAgent.toLowerCase()
    console.log(`userAgent ${userAgent}`)

    if (userAgent.indexOf('electron') > -1){
        return true
    }

    return false
}