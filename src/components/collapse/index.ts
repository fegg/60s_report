import './collapse.scss';
import { attachPropertiesToComponent } from '../slider/utils/attach-properties-to-component';
import { Collapse, CollapsePanel } from './collapse';

export type { CollapseProps, CollapsePanelProps } from './collapse';

export default attachPropertiesToComponent(Collapse, {
  Panel: CollapsePanel,
});
