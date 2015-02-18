$(window).load(function(){
function runExample() {
    "use strict";
    
    var consoleSocket = null, timer = null;
    var $consoleinp = $('input[name=consoledata]');
    var $chatinp = $('input[name=chatdata]');
    var domain = window.location.hostname;
    
    showMain();
          
    $('#consoleForm').submit(sendCommand);
    $('#chatForm').submit(sendChat);
    $('#button-editor').click(pageEditor);
    $('#button-console').click(pageConsole);
    $('#button-chat').click(pageChat);
    $('#button-dashboard').click(pageDashboard);
    $('#button-start').click(start);
    $('#button-stop').click(stop);
    $('#button-refresh').click(refresh);
    $('#editsubmit').click(updateFile);
    
    function showMain(){
        loadConsole();
        loadConfigs();
        loadDashboard();
        startTimer();
        if(window.location.href.indexOf("editor") > -1) {
           pageEditor();
        } else if(window.location.href.indexOf("console") > -1) {
           pageConsole();
        } else if(window.location.href.indexOf("chat") > -1) {
           pageChat();
        } else {
            pageDashboard();
        }
    }
    
    function loadConsole() {
        $('#console').empty();
        consoleSocket = new WebSocket("ws://"+ domain +"/sock");
        consoleSocket.onmessage = function (event) {
          newMessage(event.data);
        }
        $('#console').scrollTop($('#console').height());
        $('#chat').scrollTop($('#chat').height());
    }
    
    // create a new message in the DOM after it comes
    // in from the server (via child_added)
    function newMessage(snap) {
        var $console = $('#console');
        var $chat = $('#chat');
        var txt = snap;
        var prefix = txt.substring(17, 18);
        var prefix2 = txt.substring(17, 25)
        var time = txt.substring(2, 9);
        if (prefix == "<") {
            var msg = txt.substring(16);
            var name = txt.substring(txt.lastIndexOf("<")+1,txt.lastIndexOf(">"));
            $('<li class="collection-item flow-text" /> ').html("<span class='badge'>[" + time + "]</span>" + " <strong>" + name + "</strong>: " + msg).appendTo($chat);
        } else if (prefix2 == "[Server]") {
            msg = txt.substring(26);
            $('<li class="collection-item flow-text" /> ').html("<span class='badge'>[" + time + "]</span>" + " <strong>Server</strong>: " + msg).appendTo($chat);
        } else {
            $('<li class="collection-item flow-text" /> ').text(txt).appendTo($console);
        }
        $console.scrollTop($console.height());
        $chat.scrollTop($chat.height());
    }
    
    // page functions

    function pageEditor() {
        $('#page-editor').show();
        $('#page-editor').siblings().hide();
        
        $('#button-editor').addClass("active");
        $('#button-editor').siblings().removeClass("active");
    }
    
    function pageChat() {
        $('#page-chat').show();
        $('#page-chat').siblings().hide();
        
        $('#button-chat').addClass("active");
        $('#button-chat').siblings().removeClass("active");
    }
    
    function pageConsole() {
        $('#page-console').show();
        $('#page-console').siblings().hide();
        
        $('#button-console').addClass("active");
        $('#button-console').siblings().removeClass("active");
    }
    
    function pageDashboard() {
        $('#page-dashboard').show();
        $('#page-dashboard').siblings().hide();
        
        $('#button-dashboard').addClass("active");
        $('#button-dashboard').siblings().removeClass("active");
        
        loadDashboard();
    }
    
    function launchUpload() {
        $('#modal1').openModal();
    }
    
     // post the forms contents and attempt to write a message
    function sendCommand(e) {
        e.preventDefault();
        var val = $consoleinp.val();
        $consoleinp.val(null);
        consoleSocket.send(val);
    }
    
    function sendChat(e) {
        e.preventDefault();
        var val = $chatinp.val();
        $chatinp.val(null);
        consoleSocket.send("say " + val);
    }
    
    function loadConfigs() {
        var $files = $('#files');
        $files.empty();
        $.getJSON( "/configs", function( data ) {
          var items = [];
          $.each( data, function( key, val ) {
            $files.append("<li><div class='file-name collapsible-header flow-text truncate'><a class='file' id='file"+ key +"' href='#editor'><i class='mdi-editor-insert-drive-file'></i>"+ val +"</a></div></li>");
            $('#file' + key).click(function(){
                editFile(val, key);
            });
          });
        });
    }
    
    function startTimer()
    {
        setInterval(function(){ timerUp(); }, 1000);
    }
    
    function timerUp()
    {
        timer++;
        var resetat=30;
        if(timer == resetat){
            refresh();
        }
        var tleft=resetat-timer;
        if (tleft === 0){
            timer = 0;
        }
        document.getElementById('refreshtimer').innerHTML=tleft;
    }
    
    function start() {
        $.get( "/server/start", function( data ) {});
        toast("Server started!", 4000);
    }
    
    function stop() {
        consoleSocket.send("stop");
        toast("Server stopped!", 4000);
    }
    
    function refresh() {
        loadDashboard();
        if(window.location.href.indexOf("dashboard") > -1) {
            toast("Dashboard refreshed!", 1000);
        }
    }
    
    function loadDashboard() {
        $.getJSON( "/server", function( data ) {
          $('#playerbar').css("width", data.NumPlayers / data.MaxPlayers * 100 + "%");
          $('#players').text(data.NumPlayers + "/" + data.MaxPlayers + " players online");
          
          var mempercent = data.Memory.Used / data.Memory.Total * 100
          $('#memorybar').css("width", data.Memory.Used / data.Memory.Total * 100 + "%");
          $('#memory').text(Math.floor(mempercent) + "% of RAM used");
          
          $('#version').text(data.Version);
          $('#map').text(data.Map);
          $('#gameid').text(data.GameId);
          $('#gametype').text(data.GameType);
          $('#motd').text(data.Motd);
          
          if (data.Status === null || data.status === 0) {
            $('#statustext').text("Offline");
            $('#status').addClass("red").removeClass("green");
            $('#button-stop').addClass("disabled").removeClass("red").siblings().removeClass("disabled").addClass("white-text");
          } else {
             $('#statustext').text("Online");
             $('#status').addClass("green").removeClass("red");
             $('#button-start').addClass("disabled").removeClass("green").siblings().removeClass("disabled").addClass("white-text");
          }
        });
    }
    
    function editFile(name, id) {
        var $modal = $('#editmodal');
        $modal.openModal();
        $('#filetitle').text(name);
        $('#fileid').val(id);
        $('#filecontents').load("/config/" + id);
    }
    
    function updateFile() {
        var id = $('#fileid').val();
        var content = $('#filecontents').val();
        var newcontent = escape(content.replace(/\//g, "&#47;"));
        var url = "/update/" + id + "/" + newcontent;
        $.post( url, function( data ) {
          toast(data, 4000)
        });
    }

}

runExample();
                                                      
});//]]>  