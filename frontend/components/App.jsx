import React from 'react';
import {BrowserRouter as Router, Route, Link} from 'react-router-dom';
import {Button,Navbar,Nav,NavItem} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import HelloComponent from "./HelloComponent.jsx";
import UserRegComponent from "./UserRegComponent.jsx";
export default class App extends React.Component{
	render(){
		return (			
			<Router>
				<div>
					<Navbar collapseOnSelect>					
					<Navbar.Header>
						<Navbar.Brand>
							<a href="/">Hello World</a>
						</Navbar.Brand>
						<Navbar.Toggle/>
					</Navbar.Header>				
					<Navbar.Collapse>
						<Nav>
							<LinkContainer to="/hello">
								<NavItem>Hello!</NavItem>
							</LinkContainer>
							<LinkContainer to="/register">
								<NavItem>Register</NavItem>
							</LinkContainer>
						</Nav>
					</Navbar.Collapse>
					</Navbar>
					<Route path="/hello" component={HelloComponent}/>
					<Route path="/register" component={UserRegComponent}/>
				</div>
			</Router>
		)
	}
}