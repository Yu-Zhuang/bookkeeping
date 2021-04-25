const register = async () => {
    const data = await getData()
    const res = await fetch((hostUrl + 'api/register'), {
        method: 'POST',
        body: data,
        headers: {
            'Content-Type':'application/json'
        }
    })
    //const myJson = await res.json()
    if (res.status == 200) {
        window.location.href = hostUrl
        return
    }
    else {
        alert("申請失敗") 
        window.location.href = hostUrl + "register"            
    }
    return
}
function getData() {
    id = String(document.getElementById("id").value)
    psd = String(document.getElementById("password").value)
    cpsd = String(document.getElementById("checkPassword").value)
    _name = String(document.getElementById("name").value)
    email = String(document.getElementById("email").value)

    if (psd != cpsd) {
        alert("密碼與確認密碼不符") 
        window.location.href = hostUrl + "register"            
    }
    if (id == "" || pas == "" || _name == "" || email == "") {
        alert("填寫不完全")
        window.location.href = hostUrl + "register"  
    }
    ret = {
            "id": id,
            "password": psd,
            "name": _name,
            "email": email
    }
    return JSON.stringify(ret)
}