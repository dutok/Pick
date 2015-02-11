$(window).load(function(){
function runExample() {
    "use strict";
    
    var uid = null, token = null, messages = null, submitted = false, sub = null, members = {}, userName;
    var ref = new Firebase("https://go-mine.firebaseio.com");
    var $inp = $('input[name=data]');
    var $join = $('#joinForm');
    
    // handle input and form events
    $('#main-layer').hide();
    $('#chatForm').submit(sendCommand);
    $('#login').click(authenticate);
    $('a[href="#logout"]').click(logout);
    
    function authenticate(e) {
        e.preventDefault();
        ref.authWithOAuthPopup('github', function(err, user) {
            if (err) {
                console.log(err, 'error');
            } else if (user) {
                // logged in!
                allowSending(user);
                uid = user.uid;
                token = user.token;
                console.log('logged in with id', uid);
                $('#login-layer').hide();
                $('#main-layer').show();
                loadConsole();
            } else {
                // logged out
                $('#login-layer').show();
                $('#main-layer').hide();
            }
        },
        {remember: "default",});
    }
    
    function allowSending(user) {
        var allowed = ref.child("allowed");
        allowed.child(user.github.username).set({
          token: user.token,
        });
    }
    
    // post the forms contents and attempt to write a message
    function sendCommand(e) {
        e.preventDefault();
        submitted = true;
        var val = $inp.val();
        $inp.val(null);
        var serverurl = window.location.href.split('?')[0] + "command";
        var url = serverurl + "/" + val + "/" + token;
        $.post(url);
        console.log(url);
    }
    
    function loadConsole() {
        emptyConsole();
        messages = ref.child('console/messages').limitToLast(30);
        messages.on('child_added', newMessage);
        messages.on('child_removed', dropMessage);   
    }
    
    function emptyConsole() {
        $join.detach();
        return $('ul.chatbox').empty();
    }
    
    // create a new message in the DOM after it comes
    // in from the server (via child_added)
    function newMessage(snap) {
        var $chat = $('ul.chatbox');
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
       emptyConsole();
    }

    // print results of write attempt so we can see if
    // rules allowed or denied the attempt
    function result(err) {
        if (err) {
            log(err.code, 'error');
        } else {
            log('success!');
        }
    }

    var to;

    // clear write results after 5 seconds
    function delayedClear() {
        to && clearTimeout(to);
        to = setTimeout(clearNow, 5000);
    }

    // clear write results now
    function clearNow() {
        $('p.result').text('');
        to && clearTimeout(to);
        to = null;
        submitted = false;
    }

}

runExample();
                                                      
});//]]>  