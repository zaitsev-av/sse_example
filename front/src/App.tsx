import {useState, ChangeEvent} from 'react'
import {useGetConnectQuery, EntityType, useCreateEntityMutation} from "./api.ts";
import './App.css'
import {useSseForDocumentUpdates} from "./useEventsSSE.ts";


function App() {
	const [value, setValue] = useState<string>('')
	const [entity, setEntity] = useState<EntityType[]>([])
	const [startSse, setStartSse] = useState<boolean>(false)

	const {data} = useGetConnectQuery()
	const [create] = useCreateEntityMutation()

	const handleOnChange = (e: ChangeEvent<HTMLInputElement>) => {
		setValue(e.target.value)
	}

	useSseForDocumentUpdates(startSse, () => {
		setStartSse(false)
	})

	return (
		<>
			<div className="card">
				<input value={value} onChange={handleOnChange}/>
				<button onClick={() => {
					setEntity(prevState => [...prevState, {
						name: value,
						id: (100000 * Math.random()).toFixed().toString(),
						status: 'pending'
					} as EntityType])
					setValue("")
				}}>
					add
				</button>
				<div style={{display: "flex", justifyContent: "space-around"}}>
					<ul>
						<h3>Созданный данне</h3>
						{entity.map((el) => {
							return <li style={{
								padding: '10px', width: '200px', display: 'flex',
								flexDirection: "row",
								justifyContent: 'space-between',
							}}>{el.name} <span style={{
								background: el.status === 'pending' ? '#4d7373' : el.status === "in_progress" ? "rgb(230 194 29 / 90%)" : '#09ff066b',
								padding: '10px',
								borderRadius: '8px',
								textAlign: "left"
							}}>{el.status}</span></li>
						})}
					</ul>
					<ul>
						<h3>Полученные данне</h3>
						{data?.map((el) => {
							return <li style={{
								padding: '10px', width: '200px', display: 'flex',
								flexDirection: "row",
								justifyContent: 'space-between',
							}}>{el.name} <span style={{
								background: el.status === 'pending' ? '#4d7373' : el.status === "in_progress" ? "rgb(230 194 29 / 90%)" : '#09ff066b',
								padding: '10px',
								borderRadius: '8px',
								textAlign: "left"
							}}>{el.status}</span></li>
						})}
					</ul>
				</div>

			</div>
			<p className="read-the-docs">
				<button onClick={() => {
					create(entity).unwrap().then(()=> {
						setTimeout(()=> {
							setStartSse(true)
							setEntity([])
						}, 3000)
					})

				}}>Submit
				</button>
			</p>
		</>
	)
}

export default App
