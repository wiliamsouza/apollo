/** @jsx React.DOM */

var NewUserForm = React.createClass({
  handleSubmit: function() {
    var name = this.refs.name.getDOMNode().value.trim();
    var email = this.refs.email.getDOMNode().value.trim();
    var pwd = this.refs.pwd.getDOMNode().value.trim();
    var pass = this.refs.pass.getDOMNode().value.trim();
    // TODO: Check password is the same here!
    this.postUser({name: name, email: email, password:pass});
    this.refs.name.getDOMNode().value = '';
    this.refs.email.getDOMNode().value = '';
    this.refs.pwd.getDOMNode().value = '';
    this.refs.pass.getDOMNode().value = '';
    return false;
  },
  postUser: function(user) {
    $.ajax({
      url: this.props.url,
      type: 'POST',
      dataType: 'json',
      data: JSON.stringify(user),
      success: function(user) {
	      console.log(user)
      }.bind(this)
    });
  },
  render: function() {
    return (
      <form className="newUserForm" onSubmit={this.handleSubmit}>
        <input type="text" placeholder="Full name." ref="name" />
        <input type="text" placeholder="Your email." ref="email" />
        <input type="password" placeholder="Password." ref="pwd" />
        <input type="password" placeholder="Password again." ref="pass" />
        <input type="submit" value="Post" />
      </form>
    );
  }
});

React.renderComponent(
  <NewUserForm url="http://localhost:8000/user" />,
  document.getElementById('newUserForm')
);
