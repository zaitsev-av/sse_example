import {FC} from "react";
import {EntityType} from "../../api";
import s from './badge.module.css'

interface Props {
	entity: EntityType
}

export const Badge: FC<Props> = ({entity}) => {
	return (
		<li className={s.root}>{entity.name} <span className={s.element} style={{
			background: entity.status === 'pending' ? '#4d7373' : entity.status === "in_progress" ? "rgb(230 194 29 / 90%)" : entity.status === 'error' ? 'rgba(255,14,6,0.66)' : '#09ff066b',
		}}>{entity.status}</span>
		</li>
	)
}
