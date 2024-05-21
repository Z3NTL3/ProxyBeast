/*
dont care about js
*/

window.runtime.EventsOn("checker_end", () => {
    Fire({
        ...STD_PROPS,
        icon: "info",
        text: "Proxy checking ended"
    })

    document.getElementById("cancel").setAttribute("class", "flex w-full justify-end h-full invisible")

    document.getElementById("started").setAttribute("style", "color: red;")
    document.getElementById("started").innerText = "false"


    document.getElementById("checked").innerText = `0`
})

window.runtime.EventsOn("checker_start", () => {
    Fire({
        ...STD_PROPS,
        icon: "info",
        text: "Proxy checking started"
    })

    document.getElementById("cancel").setAttribute("class", "flex w-full justify-end h-full visible")

    document.getElementById("started").setAttribute("style", "color: green;")
    document.getElementById("started").innerText = "true"

})

window.runtime.EventsOn("current_thread", (th) => {
    document.getElementById("checked").innerText = `${th}`
})

window.runtime.EventsOn("checker_load", (load) => {
    document.getElementById("load").innerText = `${load}`
})

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
        text: "Path could not be resolved or MKDIR operation for saves dir has failed. Closing in 6 seconds, give app permissions on FS (read/write in current app folder).",
        timer: 6000
    })
    setTimeout(() => {
        window.runtime.Quit()
    }, 3000)
})
