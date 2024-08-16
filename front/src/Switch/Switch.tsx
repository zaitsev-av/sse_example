import {SwitchProps} from "@radix-ui/react-switch";
import * as SwitchRadix from '@radix-ui/react-switch';
import {FC} from "react";
import s from './styles.module.css';

export const Switch: FC<SwitchProps> = (props) => {
	return <form>
		<div style={{display: 'flex', alignItems: 'center'}}>
			<label className={s.Label} htmlFor="airplane-mode" style={{paddingRight: 15}}>
				Отключить соединение
			</label>
			<SwitchRadix.Root className={s.SwitchRoot} id="airplane-mode" checked={props.checked}
			                  onCheckedChange={props.onCheckedChange}>
				<SwitchRadix.Thumb className={s.SwitchThumb}/>
			</SwitchRadix.Root>
		</div>
	</form>
}
