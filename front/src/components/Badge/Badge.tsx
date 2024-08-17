import {FC} from "react";
import {EntityType} from "../../api/api.ts";

interface Props {
	entity: EntityType
};
export const Badge: FC<Props> = ({entity}) => {
	return (
		<li style={{
			padding: '10px', width: '200px', display: 'flex',
			flexDirection: "row",
			justifyContent: 'space-between',
		}}>{entity.name} <span style={{
			background: entity.status === 'pending' ? '#4d7373' : entity.status === "in_progress" ? "rgb(230 194 29 / 90%)" : '#09ff066b',
			padding: '10px',
			borderRadius: '8px',
			textAlign: "left"
		}}>{entity.status}</span>
		</li>
	)
}
