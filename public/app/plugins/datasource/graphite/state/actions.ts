import { GraphiteQueryEditorAngularDependencies, GraphiteSegment, GraphiteTag } from '../types';
import { createAction } from '@reduxjs/toolkit';

/**
 * List of possible actions changing the state of QueryEditor
 */

/**
 * This is used only during the transition to react. It will be removed after migrating all components.
 */
const init = createAction<GraphiteQueryEditorAngularDependencies>('init');

// Metrics & Tags
const segmentValueChanged = createAction<{ segment: GraphiteSegment; index: number }>('segment-value-changed');

// Tags
const addNewTag = createAction<{ segment: GraphiteSegment }>('add-new-tag');
const tagChanged = createAction<{ tag: GraphiteTag; index: number }>('tag-changed');
const unpause = createAction('unpause');

// Functions
const addFunction = createAction<{ name: string }>('add-function');
const removeFunction = createAction<{ funcDef: object }>('remove-function');
const moveFunction = createAction<{ funcDef: object; offset: number }>('move-function');

// Text editor
const targetChanged = createAction('target-changed');
const toggleEditorMode = createAction('toggle-editor');

export const actions = {
  init,
  segmentValueChanged,
  tagChanged,
  addNewTag,
  unpause,
  addFunction,
  removeFunction,
  moveFunction,
  targetChanged,
  toggleEditorMode,
};