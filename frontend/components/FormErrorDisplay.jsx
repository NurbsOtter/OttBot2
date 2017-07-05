import React from 'react';
import {Panel} from 'react-bootstrap';
export default class FormErrorDisplay extends React.Component{
	render(){
		if (this.props.error !== ''){
			return(
				<Panel>{this.props.error}</Panel>
			)
		}else{
			return null;
		}
	}
}