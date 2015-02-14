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
        $('#login-layer').hide();
        $('#main-layer').show();
        if(window.location.href.indexOf("editor") > -1) {
           pageEditor();
        } else {
            pageConsole();
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
        loadConfigs();
    }
    
    function pageConsole() {
        $('#page-console').show();
        $('#page-console').siblings().hide();
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
    }
    
    function loadConfigs() {
        var $files = $('#files');
        $files.empty();
        console.log(token)
        $.getJSON( "/configs/" + token, function( data ) {
          var items = [];
          $.each( data, function( key, val ) {
            $files.append("<li><div class='collapsible-header flow-text truncate'><a class='file' id='file"+ key +"' href='#editor'><i class='mdi-editor-insert-drive-file'></i>"+ val +"</a></div></li>");
            $('#file' + key).click(function(){
                editFile(val, key);
            });
          });
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
        var newcontent = escape(content);
        var url = "/update/" + id + "/" + newcontent + "/" + token;
        console.log(url);
        $.post( url, function( data ) {
          toast(data, 4000)
        });
    }

}

runExample();
                                                      
});//]]>  