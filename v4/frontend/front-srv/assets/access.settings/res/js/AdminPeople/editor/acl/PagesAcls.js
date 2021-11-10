import React from 'react'
import LangUtils from 'pydio/util/lang'
import {IdmWorkspace, TreeNode} from 'cells-sdk';
import WorkspaceAcl from './WorkspaceAcl'

class PagesAcls extends React.Component{

    constructor(props){
        super(props);
        const m = (id) => props.pydio.MessageHash['pydio_role.' + id] || id;

        let workspaces = [];
        const homepageWorkspace = new IdmWorkspace();
        homepageWorkspace.UUID = "homepage";
        homepageWorkspace.Label = m('workspace.statics.home.title');
        homepageWorkspace.Description = m('workspace.statics.home.description');
        homepageWorkspace.Slug = "homepage";
        homepageWorkspace.RootNodes = {"homepage-ROOT": TreeNode.constructFromObject({Uuid:"homepage-ROOT"})};
        workspaces.push(homepageWorkspace);
        if(props.showSettings) {
            const settingsWorkspace = new IdmWorkspace();
            settingsWorkspace.UUID = "settings";
            settingsWorkspace.Label = m('workspace.statics.settings.title');
            settingsWorkspace.Description = m('workspace.statics.settings.description');
            settingsWorkspace.Slug = "settings";
            settingsWorkspace.RootNodes = {"settings-ROOT": TreeNode.constructFromObject({Uuid:"settings-ROOT"})};
            workspaces.push(settingsWorkspace);
        }
        workspaces.sort(LangUtils.arraySorter('Label', false, true));
        this.state = {workspaces};
    }

    render(){
        const {role} = this.props;
        const {workspaces} = this.state;
        if(!role){
            return <div></div>
        }
        return (
            <div style={{backgroundColor:'white'}} className={"material-list"}>{workspaces.map( ws => <WorkspaceAcl workspace={ws} role={role} /> )}</div>
        );

    }

}

export {PagesAcls as default}