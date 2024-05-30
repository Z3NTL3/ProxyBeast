/*
dont care about js
*/

setInterval(async () => {
    console.log('Alltime goods: ', await window.go.core.App.GetAllTimeGoods())
    console.log('Alltime scans: ', await window.go.core.App.GetAllTimeScans())
}, 2000)
window.runtime.EventsOn("checker_end", async () => {
    Fire({
        ...STD_PROPS,
        icon: "info",
        text: "Proxy checking ended"
    })

    document.getElementById("cancel").setAttribute("class", "flex w-full justify-end h-full invisible")
    document.getElementById("TotalScans").innerText = await window.go.core.App.GetAllTimeScans()
})

window.runtime.EventsOn("checker_start", () => {
    Fire({
        ...STD_PROPS,
        icon: "info",
        text: "Proxy checking started"
    })

    document.getElementById("cancel").setAttribute("class", "flex items-center ml-[15px] text-white font-bold")
    document.getElementById("results-table").innerHTML = "Starting... Please wait, setting up workers."
})

window.runtime.EventsOn("proxy_data", async (data) => {
    data = JSON.parse(data)
    let current_data = document.getElementById("results-table").innerHTML
    let apndData = `
        <tr class="flex w-full text-gray-400 bg-base-300 rounded-md mb-[3px]" >
            <td class="p-2 w-1/4 text-xs break-words">${data.proxy}</td>
            <td class="p-2 w-1/4 text-xs break-words">${data.protocol}</td>
            <td class="p-2 w-1/4 text-xs break-words">${data.latency}</td>
            <td class="p-2 w-1/4 text-xs break-words">${data.anonimity}</td>
        </tr>
    `

    document.getElementById("results-table").innerHTML = apndData + current_data
    document.getElementById("TotalGoods").innerText = await window.go.core.App.GetAllTimeGoods()
})

// window.runtime.EventsOn("current_thread", (th) => {
//     document.getElementById("checked").innerText = `${th}`
// })

// window.runtime.EventsOn("checker_load", (load) => {
//     document.getElementById("load").innerText = `${load}`
// })

window.runtime.EventsOn("proto_unknown", () => Fire({
    ...STD_PROPS,
    icon: "error",
    text: "Protocol can only be one of: HTTP/HTTPS/SOCKS4/SOCKS5"
}))
window.runtime.EventsOn("dialog_save_file", function (loc) {
    console.log("got2")
    Fire({
        ...STD_PROPS,
        icon: "success",
        text: "Synced to save location"
    })
    document.querySelectorAll("#inputs")[1].setAttribute("placeholder", loc)
})

window.runtime.EventsOn("dialog_input_file", function (loc) {
    console.log("got1")
    Fire({
        ...STD_PROPS,
        icon: "success",
        text: "Proxy list has been loaded"
    })
    document.querySelector("#inputs").setAttribute("placeholder", loc)
})

window.runtime.EventsOn("msg", function (payload) {
    Fire({
        ...STD_PROPS,
        icon: "info",
        text: payload
    })
})

window.runtime.EventsOn("error", function (err) {
    Fire({
        ...STD_PROPS,
        icon: "error",
        text: err
    })
})

window.runtime.EventsOn("svdir_failure", function (_) {
    Fire({
        ...STD_PROPS,
        icon: "warning",
        title: "Warning",
        text: "Failing app permissions, cannot access filesystem. Try to open me as Administrator. Exiting automatically in 6 seconds.",
        timer: 6000
    })
    setTimeout(() => {
        window.runtime.Quit()
    }, 3000)
})
