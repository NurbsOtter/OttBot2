import React from 'react';
import {Form,FormGroup,ControlLabel,FormControl,Button,Panel} from 'react-bootstrap';
import {Router} from 'react-router';
import FormErrorDisplay from './FormErrorDisplay.jsx';
export default class UserRegComponent extends React.Component{
	constructor(props){
		super(props);
		this.state= {
			userName:'',
			email:'',
			password:'',
			error:''
		};
		this.handleUsrChange = this.handleUsrChange.bind(this);
		this.handleEmailChange = this.handleEmailChange.bind(this);
		this.handlePasswordChange = this.handlePasswordChange.bind(this);
		this.doRegister = this.doRegister.bind(this);
		this.validateUserName = this.validateUserName.bind(this);
	}
	handleUsrChange(e){
		this.setState({userName:e.target.value});
	}
	validateUserName(){
		if (this.state.userName != ""){
			return 'success';
		}else{
			return 'error';
		}
	}
	handleEmailChange(e){
		this.setState({email:e.target.value});
	}
	handlePasswordChange(e){
		this.setState({password:e.target.value});
	}
	doRegister(){
		fetch('/user/register',{
			method:'POST',
			body:JSON.stringify({
				'UserName':this.state.userName,
				'Email':this.state.email,
				'password':this.state.password
			})
		}).then((response)=>{
			if (response.ok){
				return response.json();
			}
			throw new Error("Username or Password taken.");
		}).then((data)=>{
			if (data.success == 0){
				this.setState({error:data.message});
			}else{
				this.props.history.push('/');
			}			
		}).catch((error)=>{
			this.setState({error:error.message});
		})
	}
	render(){
		return(
			<div>
				<Form>
					<FormGroup>
						<FormErrorDisplay error={this.state.error}/>
					</FormGroup>
					<FormGroup controlId = "formUserName" validationState={this.validateUserName()}>
						<ControlLabel>User Name</ControlLabel>
						<FormControl
							type="text"
							value={this.state.userName}
							placeholder="Enter UserName"
							onChange={this.handleUsrChange}
						/>
					</FormGroup>
					<FormGroup controlId = "formEmail">
						<ControlLabel>Email Address</ControlLabel>
						<FormControl
							type="text"
							value={this.state.email}
							placeholder="Enter Email"
							onChange={this.handleEmailChange}
						/>
					</FormGroup>
					<FormGroup controlId = "formPassword">
						<ControlLabel>Password</ControlLabel>
						<FormControl
							type="password"
							value={this.state.password}
							placeholder="Enter Password"
							onChange={this.handlePasswordChange}
						/>
					</FormGroup>
					<FormGroup controlId = "submitButton">
						<Button bsStyle="primary" onClick={this.doRegister}>Register</Button>
					</FormGroup>
				</Form>
				<h1>{this.state.userName}</h1>
			</div>
		)
	}
}