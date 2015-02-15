$(window).load(function(){
function runExample() {
    "use strict";
    
    var uid = null, email = null, username = null, avatar = null, token = null, messages = null, sub = null, members = {}, userName;
    var ref = new Firebase("https://go-mine.firebaseio.com");
    var $inp = $('input[name=data]');
    var $join = $('#joinForm');
    
    // handle input and form events
    var authData = ref.getAuth();
    if (authData) {
        setUser(authData);
        showMain();
    } else {
        showLogin();
    }
          
    $('#chatForm').submit(sendCommand);
    $('#login').click(authenticate);
    $('a[href="#logout"]').click(logout);
    $('#button-editor').click(pageEditor);
    $('#button-console').click(pageConsole);
    $('#button-dashboard').click(pageDashboard);
    $('#button-start').click(start);
    $('#button-stop').click(stop);
    $('#editsubmit').click(updateFile);
    
    function authenticate(e) {
        e.preventDefault();
        ref.authWithOAuthPopup('github', function(err, user) {
            if (err) {
                console.log(err, 'error');
            } else if (user) {
                // logged in!
                setUser(user);
                showMain();
            } else {
                // logged out
                showLogin();
            }
        },
        {remember: "default",});
    }
    
    function setUser(user){
        allowSending(user);
        uid = user.uid;
        token = user.github.accessToken;
        avatar = user.github.cachedUserProfile.avatar_url;
        username = user.github.username;
        email = user.github.email;
    }
    
    function showMain(){
        loadConsole();
        loadConfigs();
        loadDashboard();
        $('#login-layer').hide();
        $('#main-layer').show();
        if(window.location.href.indexOf("editor") > -1) {
           pageEditor();
        } else if(window.location.href.indexOf("console") > -1) {
           pageConsole();
        } else {
            pageDashboard();
        }
    }
    
    function showLogin() {
        $('#login-layer').show();
        $('#main-layer').hide();
    }
    
    function allowSending(user) {
        var allowed = ref.child("allowed");
        allowed.child(user.github.accessToken).set({
          name: user.github.username,
        });
    }
    
    function loadConsole() {
        $('#chatbox').empty();
        messages = ref.child('console/messages').limitToLast(30);
        messages.on('child_added', newMessage);
        messages.on('child_removed', dropMessage);   
    }
    
    // create a new message in the DOM after it comes
    // in from the server (via child_added)
    function newMessage(snap) {
        var $chat = $('#chatbox');
        var dat = snap.val();
        var txt = dat.Body;
        $('<li class="collection-item flow-text" /> ').attr('data-id', snap.key()).text(txt).appendTo($chat);
        $chat.scrollTop($chat.height());
    }
    
    // remove message locally after child_removed
    function dropMessage(snap) {
        $('li[data-id="'+snap.key()+'"]').remove();
    }
    
    function logout(e) {
       e.preventDefault();
       ref.unauth();
       $('#login-layer').show();
       $('#main-layer').hide();
    }
    
    // page functions

    function pageEditor() {
        $('#page-editor').show();
        $('#page-editor').siblings().hide();
        
        $('#button-editor').addClass("active");
        $('#button-editor').siblings().removeClass("active");
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
        var val = $inp.val();
        $inp.val(null);
        var serverurl = "http://" + window.location.host + "/command";
        var url = serverurl + "/" + val + "/" + token;
        $.post(url);
        console.log(url)
    }
    
    function loadConfigs() {
        var $files = $('#files');
        $files.empty();
        $.getJSON( "/configs/" + token, function( data ) {
          var items = [];
          $.each( data, function( key, val ) {
            $files.append("<li><div class='file-name collapsible-header flow-text truncate'><a class='file' id='file"+ key +"' href='#editor'><i class='mdi-editor-insert-drive-file'></i>"+ val +"</a></div></li>");
            $('#file' + key).click(function(){
                editFile(val, key);
            });
          });
        });
    }
    
    function start() {
        $.get( "/server/start", function( data ) {});
        toast("Server started!", 4000);
    }
    
    function stop() {
        $.get( "/server/stop", function( data ) {});
        toast("Server stopped!", 4000);
    }
    
    function loadDashboard() {
        $.getJSON( "/server", function( data ) {
          console.log(data);
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
          } else {
             $('#statustext').text("Online");
             $('#status').addClass("green").removeClass("re");
          }
        });
    }
    
    function editFile(name, id) {
        var $modal = $('#editmodal');
        $modal.openModal();
        $('#filetitle').text(name);
        $('#fileid').val(id);
        $('#filecontents').load("/config/" + id + "/" + token);
    }
    
    function updateFile() {
        console.log("Updated!");
        var id = $('#fileid').val();
        var content = $('#filecontents').val();
        var newcontent = escape(content.replace(/\//g, "&#47;"));
        var url = "/update/" + id + "/" + newcontent + "/" + token;
        console.log(url);
        $.post( url, function( data ) {
          toast(data, 4000)
        });
    }

}

runExample();
                                                      
});//]]>  