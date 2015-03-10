$(window).load(function(){
function runExample() {
    "use strict";
    
    var consoleSocket = null, timer = null, starttime = null;
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
    $('#accountName').click(loadAccount);
    
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
        var $console = $('#console');
        var $chat = $('#chat');
        $console.empty();
        consoleSocket = new WebSocket("ws://"+ domain +"/sock");
        consoleSocket.onmessage = function (event) {
          newMessage(event.data);
        }
        $console.scrollTop($console[0].scrollHeight);
        $chat.scrollTop($chat[0].scrollHeight);
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
            $('<li class="collection-item flow-text" /> ').html("<span class='badge'>[" + time + "]</span>" + " <strong class='server-chat'>Server</strong>: " + msg).appendTo($chat);
        } else {
            $('<li class="collection-item flow-text" /> ').text(txt).appendTo($console);
        }
        $console.scrollTop($console[0].scrollHeight);
        $chat.scrollTop($chat[0].scrollHeight);
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
        
        $('#console').scrollTop($('#console')[0].scrollHeight);
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
    
    function refresh() {
        loadDashboard();
        timer = 0;
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
            loadDashboard();
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
    
    function loadDashboard() {
        $.getJSON( "/server", function( data ) {
          $('#playerbar').css("width", data.NumPlayers / data.MaxPlayers * 100 + "%");
          $('#players').text(data.NumPlayers + "/" + data.MaxPlayers + " players online");
          
          var mempercent = data.Memory.Used / data.Memory.Total * 100
          $('#memorybar').css("width", data.Memory.Used / data.Memory.Total * 100 + "%");
          $('#memory').text(Math.floor(mempercent) + "% of RAM used");
          
          $('#cpubar').css("width", Math.floor(data.CPU) + "%");
          $('#cpu').text(Math.floor(data.CPU) + "% of CPU used");
          
          $('#tpsbar').css("width", Math.floor(data.Tps) / 20 * 100 + "%");
          $('#tps').text(data.Tps + " ticks per second (TPS)");
          
          $('#version').text(data.Version);
          $('#map').text(data.Map);
          $('#gameid').text(data.GameId);
          $('#gametype').text(data.GameType);
          $('#motd').text(data.Motd);
          
          if (data.Status === null || data.status === 0) {
              $('#uptime').empty();
          } else {
              if (starttime != data.StartTime) {
                starttime = data.StartTime;
                $('#uptime').html("Started <span data-livestamp='"+ data.StartTime +"'></span>")
              }
          }
          
          if (data.Status === null || data.status === 0) {
            $('#statustext').text("Offline");
            $('#status').css("color", "#f44336").addClass("mdi-content-clear").removeClass("mdi-navigation-check");
          } else {
             $('#statustext').text("Online");
             $('#status').css("color", "#4CAF50").addClass("mdi-navigation-check").removeClass("mdi-content-clear");
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
    
    function loadAccount() {
        var $modal = $('#accountmodal')
        $modal.openModal();
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