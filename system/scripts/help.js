const messageLifetime = 600 * 1000;
const maxNode=15;
const minNode=1;

var currentNodes = {}
var activeNodes = {}
var messages = {}

var Now = function () {
    return (new Date()).getTime()
}
function calcNode(){
    let channles=Object.keys(activeNodes)
    currentNodes = {}
    for (var i=0;i<channles.length;i++){
        let channel=channles[i]
        let nodeids=Object.keys(activeNodes[channel])
        let list=[]
        for(var nodeidindex=0;nodeidindex<nodeids.length;nodeidindex++){
            list.push({
                nodeid:nodeids[nodeidindex],
                count:activeNodes[channel][nodeids[nodeidindex]]
            })
        }
        list.sort(function(a,b){return b.count-a.count})
        let result=list.slice(0,maxNode)
        currentNodes[channel]=[]
        for(var ri=0;ri<result.length;ri++){
            currentNodes[channel].push(result[ri].nodeid)
        }
    }
    activeNodes={}
}
function gc() {
    let keys = Object.keys(messages)
    let now = Now()
    for (var i = 0; i < keys.length; i++) {
        if (now - messages[keys[i]].Timestamp > messageLifetime) {
            delete messages[keys[i]]
        }
    }
}
function reply(msgid, msg) {
    let data = messages[msgid]
    if (data) {
        delete messages[msgid]
        Send(data.ID, msg)
    }
}
function onActive(id, channel) {
    if (activeNodes[channel] == null) {
        activeNodes[channel] = {}
    }
    if (activeNodes[channel][id] == null) {
        activeNodes[channel][id] = 0
    }
    activeNodes[channel][id]++
}

function newMessage(msgid, id) {
    messages[msgid] = {
        ID: id,
        Timestamp: Now()
    }
}
function OnTicker() {
    gc()
    calcNode()
}
function broadcast(channel, msg) {
    if (currentNodes[channel]&&currentNodes[channel].length>=minNode) {
        for (var i=0;i<currentNodes[channel].length;i++){
            Send(currentNodes[channel][i],msg)
        }
        return false
    }
    return true
}
function OnMessage(id, msg) {
    let data = msg.split(" ")
    if (data.length > 2) {
        if (data[1] == "help") {
            onActive(id, data[2])
            newMessage(data[2] + " " + data[3], id)
            return broadcast(data[2],msg)
        } else if (data[1] == 'found') {
            onActive(id, data[2])
            reply(data[2] + " " + data[3].split("|")[0], msg)
            return false
        }
    }
    return true
}