/*
 * Copyright 2007-2019 Charles du Jeu - Abstrium SAS <team (at) pyd.io>
 * This file is part of Pydio.
 *
 * Pydio is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */
import React from 'react'
import {TextField, SelectField, AutoComplete} from 'material-ui'

const noWrap = {
    whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis'
};

const v2Block = {
    backgroundColor:'rgb(246, 246, 248)',
    borderRadius:'3px 3px 0 0',
    height:52,
    marginTop: 8
}

const styles = {
    textField:{
        inputStyle:{backgroundColor:'rgba(224, 224, 224, 0.33)',height: 34, borderRadius: 3, marginTop: 6, padding: 7},
        hintStyle:{paddingLeft: 7, color:'rgba(0,0,0,0.5)', ...noWrap, width: '100%'},
        underlineStyle:{opacity:0},
        underlineFocusStyle:{opacity:1, borderRadius: '0px 0px 3px 3px'},
        errorStyle:{bottom:-4}
    },
    textFieldV2:{
        style:{...v2Block},
        inputStyle:{position: 'absolute', height:30, marginTop:0, bottom: 2, paddingLeft: 8, paddingRight: 8},
        hintStyle:{bottom: 4, paddingLeft: 7, color:'rgba(0,0,0,0.5)', ...noWrap, width: '100%'},
        underlineStyle:{opacity:1, bottom: 0},
        underlineFocusStyle:{opacity:1, borderRadius: 0, bottom: 0},
        floatingLabelFixed: true,
        floatingLabelStyle:{top:26, left: 8, width:'127%', ...noWrap},
        floatingLabelShrinkStyle:{top:26, left: 8},
        errorStyle:{position:'absolute', bottom:8, right:8}
    },
    textareaField:{
        rows: 4,
        rowsMax: 4,
        inputStyle:{backgroundColor:'rgba(224, 224, 224, 0.33)',height: 106, borderRadius: 3, marginTop: 6, padding: 7},
        textareaStyle:{marginTop: 0, marginBottom: 0},
        hintStyle:{paddingLeft: 7, color:'rgba(0,0,0,0.5)', ...noWrap, width: '100%', top: 12, bottom: 'inherit'},
        underlineStyle:{opacity:0},
        underlineFocusStyle:{opacity:1, borderRadius: '0px 0px 3px 3px'},
        errorStyle:{bottom: -3}
    },
    textareaFieldV2:{
        rows: 4,
        rowsMax: 4,
        style:{height: 128},
        inputStyle:{backgroundColor:v2Block.backgroundColor, height: 120, borderRadius: v2Block.borderRadius, marginTop: 8, paddingLeft: 8},
        textareaStyle:{marginTop: 24, marginBottom: 0},
        floatingLabelFixed: true,
        floatingLabelStyle:{top:35, left:6, width:'127%', ...noWrap},
        floatingLabelShrinkStyle:{top:35, left: 6},
        hintStyle:{paddingLeft: 7, color:'rgba(0,0,0,0.5)', ...noWrap, width: '100%', top: 12, bottom: 'inherit'},
        underlineStyle:{opacity:1, bottom: 0},
        underlineFocusStyle:{opacity:1, bottom: 0, borderRadius: '0px 0px 3px 3px'},
        errorStyle:{position:'absolute', bottom:8, right:8}
    },
    selectField:{
        style:{backgroundColor:'rgba(224, 224, 224, 0.33)',height: 34, borderRadius: 3, marginTop: 6, padding: 7, paddingRight: 0, overflow:'hidden'},
        menuStyle:{marginTop: -12},
        hintStyle:{paddingLeft: 0, marginBottom: -7, paddingRight:56, color:'rgba(0,0,0,0.34)', ...noWrap, width: '100%'},
        underlineShow: false
    },
    selectFieldV2:{
        style:{...v2Block, padding: 8, paddingRight: 0, overflow:'hidden'},
        menuStyle:{marginTop: -6},
        hintStyle:{paddingLeft: 0, marginBottom: -7, paddingRight:56, color:'rgba(0,0,0,0.34)', ...noWrap, width: '100%'},
        underlineStyle: {opacity:1, bottom: 0, left: 0, right: 0},
        underlineFocusStyle:{opacity:1, borderRadius: 0, bottom: 0},
        floatingLabelFixed: true,
        floatingLabelStyle:{top:26, left: 8, width:'127%', ...noWrap},
        floatingLabelShrinkStyle:{top:26, left: 8},
        dropDownMenuProps:{iconStyle:{right: 0, fill: '#9e9e9e'}}
    },
    div:{
        backgroundColor:'rgba(224, 224, 224, 0.33)', color:'rgba(0,0,0,.5)',
        height: 34, borderRadius: 3, marginTop: 6, padding: 7, paddingRight: 0
    },
    toggleField:{
        style: {
            backgroundColor: 'rgba(224, 224, 224, 0.33)',
            padding: '7px 5px 4px',
            borderRadius: 3,
            fontSize: 15,
            margin:'6px 0 7px'
        }
    },
    toggleFieldV2:{
        style:{
            ...v2Block,
            borderRadius: 4,
            fontSize: 15,
            /*
            backgroundColor:'transparent',
            border: '1px solid rgb(224, 224, 224)',
             */
            padding: '15px 10px 4px'
        }
    },
    fillBlockV2Right:{
        ...v2Block,
        borderRadius:'0 4px 0 0',
        borderBottom:'1px solid rgb(224, 224, 224)'
    },
    fillBlockV2Left:{
        ...v2Block,
        borderRadius:'4px 0 0 0',
        borderBottom:'1px solid rgb(224, 224, 224)'
    }
};

function getV2WithBlocks(styles, hasLeft, hasRight){
    if(styles.style){
        styles.style.borderRadius = (hasLeft?'0 ':'4px ') + (hasRight?'0 ':'4px ') + '0 0';
    }
    return styles;
}

function withModernTheme(formComponent) {

    class ModernThemeComponent extends React.Component {

        mergedProps(styleProps){
            const props = this.props;
            Object.keys(props).forEach((k) => {
                if(styleProps[k]){
                    styleProps[k] = {...styleProps[k], ...props[k]};
                }
            });
            return styleProps;
        }

        componentDidMount(){
            if (this.props.focusOnMount && this.refs.component){
                this.refs.component.focus();
            }
        }

        focus(){
            if(this.refs.component){
                this.refs.component.focus();
            }
        }

        getInput(){
            if(this.refs.component){
                return this.refs.component.input;
            }
        }

        getValue(){
            return this.refs.component.getValue();
        }

        render() {

            let {variant, hasLeftBlock, hasRightBlock, ...otherProps} = this.props;
            if(variant === 'v2' || formComponent === AutoComplete) {
                if(!otherProps.floatingLabelText){
                    otherProps.floatingLabelText = otherProps.hintText
                    delete(otherProps.hintText);
                }
            } else {
                if(otherProps.floatingLabelText){
                    otherProps.hintText = otherProps.floatingLabelText;
                    delete(otherProps.floatingLabelText);
                }
            }

            if (formComponent === TextField) {
                let styleProps;
                if(this.props.multiLine){
                    if(variant === 'v2') {
                        styleProps = this.mergedProps({...styles.textareaFieldV2});
                    } else {
                        styleProps = this.mergedProps({...styles.textareaField});
                    }
                } else {
                    if(variant === 'v2') {
                        styleProps = this.mergedProps(getV2WithBlocks({...styles.textFieldV2}, hasLeftBlock, hasRightBlock));
                    } else {
                        styleProps = this.mergedProps({...styles.textField});
                    }
                }
                return <TextField {...otherProps} {...styleProps} ref={"component"} />
            } else if (formComponent === SelectField) {
                let styleProps;
                if (variant === 'v2') {
                    styleProps = this.mergedProps(getV2WithBlocks({...styles.selectFieldV2}, hasLeftBlock, hasRightBlock));
                } else {
                    styleProps = this.mergedProps({...styles.selectField});
                }
                return <SelectField {...otherProps} {...styleProps} ref={"component"}/>
            } else if (formComponent === AutoComplete) {

                const {style, ...tfStyles} = getV2WithBlocks({...styles.textFieldV2}, hasLeftBlock, hasRightBlock)
                return <AutoComplete
                    {...otherProps}
                    ref={"component"}
                    textFieldStyle={style}
                    {...tfStyles}
                />

            } else {
                return formComponent;
            }
        }
    }

    return ModernThemeComponent;

}

const ModernTextField = withModernTheme(TextField);
const ModernSelectField = withModernTheme(SelectField);
const ModernAutoComplete = withModernTheme(AutoComplete);
export {ModernTextField, ModernSelectField, ModernAutoComplete, withModernTheme, styles as ModernStyles}