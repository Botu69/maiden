import { connect } from 'react-redux';
import { createSelector } from 'reselect';
import EditActivity from './edit-activity';
import { MATRON_COMPONENT } from './constants';
import { nodeForResource } from './model/listing';

import {
    rootList,
    bufferRead,
    directoryRead,
    bufferSave,
    bufferChange,
    bufferSelect,
    scriptNew,
    scriptDuplicate,
    resourceDelete,
    resourceRename,

    toolInvoke,

    explorerActiveNode,
    explorerToggleNode,
} from './model/edit-actions';

import {
    sidebarToggle,
    sidebarSize,
    replToggle,
    replSize,
} from './model/ui-actions';

import {
    replSend,
} from './model/repl-actions';

const getBuffers = (editState) => editState.buffers;
const getActiveBuffer = (editState) => editState.activeBuffer;
const getRootNodes = (editState) => editState.rootNodes;
const getExpandedNodes = (editState) => editState.expandedNodes;

const getExplorerData = createSelector(
    [getBuffers, getActiveBuffer, getRootNodes, getExpandedNodes],
    (buffers, activeBuffer, rootNodes, expandedNodes) => {
    // enrich script listing w/ modification state, etc.

    let enrich = (items) => {
        return items.map(l => {
            let item = {...l}
            item.active = l.url === activeBuffer;
            item.toggled = expandedNodes.has(l.url);

            let buffer = buffers.get(l.url);
            if (buffer) {
                item.loaded = true;
                item.modified = buffer.get('modified') || false;
            }

            if (item.children) {
                item.children = enrich(item.children)
            }

            return item;
        })
    };

    return enrich(rootNodes.toJS())
});

const getActiveNode = createSelector(
    [getActiveBuffer, getRootNodes],
    (activeBuffer, rootNodes) => {
        return nodeForResource(rootNodes, activeBuffer)
    }
)

const mapStateToProps = (state) => {
    let {activeBuffer, buffers} = state.edit;
    return {
        activeBuffer,
        activeNode: getActiveNode(state.edit),
        buffers,
        ui: state.ui,
        explorerData: getExplorerData(state.edit),
    }
}

const mapDispatchToProps = (dispatch) => {
    return {
        // scripts
        scriptList: (api) => {
            dispatch(rootList('scripts', api))
        },
        dataList: (api) => {
            dispatch(rootList('data', api))
        },
        audioList: (api) => {
            dispatch(rootList('audio', api))
        },
        bufferRead: (api, resource) => {
            dispatch(bufferRead(api, resource))
        },
        directoryRead: (api, resource) => {
            dispatch(directoryRead(api, resource))
        },
        bufferChange: (resource, value) => {
            dispatch(bufferChange(resource, value))
        },
        bufferSelect: (resource) => {
            dispatch(bufferSelect(resource))
        },
        bufferSave: (api, resource, code, completionCB = () => {}) => {
            dispatch(bufferSave(api, resource, code, completionCB))
        },
        scriptRun: (api, resource) => {
            let file = api.fileFromResource(resource)
            let cmd = `norns.script.load("${file}")`
            dispatch(replSend(MATRON_COMPONENT, cmd))
        },

        // ui
        sidebarToggle: () => {
            dispatch(sidebarToggle())
        },
        sidebarSize: (width) => {
            dispatch(sidebarSize(width))
        },
        replToggle: () => {
            dispatch(replToggle())
        },
        replSize: (height) => {
            dispatch(replSize(height))
        },


        // tools
        toolInvoke: (name) => {
            dispatch(toolInvoke(name))
        },

        // explorer
        explorerActiveNode: (node) => {
            dispatch(explorerActiveNode(node))
        },
        explorerToggleNode: (node, toggled) => {
            dispatch(explorerToggleNode(node, toggled))
        },
        explorerScriptNew: (sibling, value) => {
            dispatch(scriptNew(sibling, value))
        },
        explorerScriptDuplicate: (source) => {
            dispatch(scriptDuplicate(source))
        },
        explorerResourceDelete: (api, resource) => {
            dispatch(resourceDelete(api, resource))
        },
        explorerResourceRename: (api, activeNode, newName, virtual) => {
            dispatch(resourceRename(api, activeNode, newName, virtual))
        },
    }
}

const BoundEditActivity = connect(
    mapStateToProps,
    mapDispatchToProps
)(EditActivity);

export default BoundEditActivity;