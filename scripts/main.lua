local myModule = require("MyModule")
require("scripts/other")

myModule.sayHi("hello lua")

function Reply(msg)
    myModule.sayHi("lua function [Reply] say: " .. msg)
end