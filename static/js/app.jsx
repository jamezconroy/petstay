var App = React.createClass({
  componentWillMount: function() {
    this.setupAjax();
    this.parseHash();
    this.setState();
  },
  setupAjax: function() {
    $.ajaxSetup({
      'beforeSend': function(xhr) {
        if (localStorage.getItem('access_token')) {
          xhr.setRequestHeader('Authorization',
                'Bearer ' + localStorage.getItem('access_token'));
        }
      }
    });
  },
  parseHash: function(){
    this.auth0 = new auth0.WebAuth({
      domain:       AUTH0_DOMAIN,
      clientID:     AUTH0_CLIENT_ID
    });
    this.auth0.parseHash(window.location.hash, function(err, authResult) {
      if (err) {
        return console.log(err);
      }
      console.log(authResult);
      if(authResult !== null && authResult.accessToken !== null && authResult.idToken !== null){
        localStorage.setItem('access_token', authResult.accessToken);
        localStorage.setItem('id_token', authResult.idToken);
        localStorage.setItem('profile', JSON.stringify(authResult.idTokenPayload));
      }
    });
  },
  setState: function(){
    var idToken = localStorage.getItem('id_token');
    if(idToken){
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  },
  render: function() {
    
    if (this.loggedIn) {
      return (<LoggedIn />);
    } else {
      return (<Home />);
    }
  }
});

var Home = React.createClass({
  authenticate: function(){
    this.webAuth = new auth0.WebAuth({
      domain:       AUTH0_DOMAIN,
      clientID:     AUTH0_CLIENT_ID,
      scope:        'openid profile',
      audience:     AUTH0_API_AUDIENCE,
      responseType: 'token id_token',
      redirectUri : AUTH0_CALLBACK_URL
    });
    this.webAuth.authorize();
  },
  render: function() {
    return (
    <div className="container">
      <div className="col-xs-12 jumbotron text-center">
        <h1>Petstay</h1>
        <p>Pets list.</p>
        <a onClick={this.authenticate} className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
      </div>
    </div>);
  }
});

var LoggedIn = React.createClass({
  logout : function(){
    localStorage.removeItem('id_token');
    localStorage.removeItem('access_token');
    localStorage.removeItem('profile');
  },

  getInitialState: function() {
    return {
      pets: []
    }
  },
  componentDidMount: function() {
    this.serverRequest = $.get('http://127.0.0.1:3000/pets?limit=3', function (result) {
      this.setState({
        pets: result,
      });
    }.bind(this));
  },

  render: function() {
    return (
      <div className="col-lg-12">
        <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
        <h2>Welcome to Avoca Petstay</h2>
        <p>Below you'll find details on all our pet visitors.</p>
        <div className="row">
        <AddPet/>
        {this.state.pets.map(function(pet, i){
          return <Pet key={i} pet={pet} />
        })}

        </div>
      </div>);
  }
});


var AddPet = React.createClass({
    addPet: function () {
        var pet = this.props.pet;
        //{"name":"Shaun-the-sheep","owner":"Jeff"}
        this.serverRequest = $.post('http://127.0.0.1:3000/pet', JSON.stringify({Name:"Shaun-the-sheep",Owner:"Jeff"}), function (result) {
            this.setState({voted: "Upvoted"})
        }.bind(this));
    },
    render: function () {
        return (
            <div className="col-xs-4">
              <div className="panel panel-default">
                {/*
                <div className="panel-heading">{this.props.pet.name} <span
                    className="pull-right">{this.state.voted}</span></div>
                <div className="panel-body">
                    {this.props.pet.owner}
                </div>
*/}                <div className="panel-footer">
                  <a onClick={this.addPet} className="btn btn-default">
                    <span className="glyphicon glyphicon-thumbs-up"></span>
                  </a>
                  <a onClick={this.downvote} className="btn btn-default pull-right">
                    <span className="glyphicon glyphicon-thumbs-down"></span>
                  </a>
                </div>
              </div>
            </div>);
    }
})

var Pet = React.createClass({
  upvote : function(){
    var pet = this.props.pet;
    this.serverRequest = $.post('http://localhost:3000/pet/' + pet.id + '/feedback', {vote : 1}, function (result) {
      this.setState({voted: "Upvoted"})
    }.bind(this));
  },
  downvote: function(){
    var product = this.props.product;
    this.serverRequest = $.post('http://localhost:3000/pet/' + pet.id + '/feedback', {vote : -1}, function (result) {
      this.setState({voted: "Downvoted"})
    }.bind(this));
  },
  delete: function(){
    var pet = this.props.pet;
      //curl -i -X DELETE http://127.0.0.1:3000/pet/3
    $.ajax({
        type: "DELETE",
        url: 'http://127.0.0.1:3000/pet/' + pet.id,
        success: function(response) {
            console.log("successfully deleted");
        },
        error: function () {
            console.log("error");
        }

    })
      // TODO - how to signal that parent re-retrieval needed to pick up fact that delete happened
    //this.serverRequest = $.get('http://127.0.0.1:3000/pets?limit=20', function (result) {
    //    this.setState({pets: result})
    //}.bind(this));
  },
  retrieveAll: function(){
    var pet = this.props.pet;
    this.serverRequest = $.get('http://127.0.0.1:3000/pets?limit=20', function (result) {
        this.setState({pets: result})
    }.bind(this));
  },
  getInitialState: function() {
    return {
      voted: null
    }
  },
  render : function(){
    return(
    <div className="col-xs-4">
      <div className="panel panel-default">
        <div className="panel-heading">{this.props.pet.name} <span className="pull-right">{this.state.voted}</span></div>
        <div className="panel-body">
          {this.props.pet.owner}
        </div>
        <div className="panel-footer">
          <a onClick={this.upvote} className="btn btn-default">
            <span className="glyphicon glyphicon-thumbs-up"></span>
          </a>
          <a onClick={this.delete} className="btn btn-default pull-right">
            <span className="glyphicon glyphicon-thumbs-down"></span>
          </a>
        </div>
      </div>
    </div>);
  }
})

ReactDOM.render(<App />,
  document.getElementById('app'));
