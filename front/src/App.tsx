import {useState, ChangeEvent, useCallback} from 'react'
import {useGetConnectQuery, EntityType, useCreateEntityMutation, useLazyResetConnectQuery} from "./api";
import {Badge} from "./components/Badge";
import {useSseForDocumentUpdates} from "./hooks";
import './App.css'


function App() {
	const [value, setValue] = useState<string>('')
	const [entity, setEntity] = useState<EntityType[]>([])
	const [startSse, setStartSse] = useState<boolean>(false)

	const {data} = useGetConnectQuery()
	const [create] = useCreateEntityMutation()
	const [reset] = useLazyResetConnectQuery()

	const handleOnChange = (e: ChangeEvent<HTMLInputElement>) => {
		setValue(e.target.value)
	}

	const handler = useCallback(() => {
		setStartSse(false)
	},[])

	useSseForDocumentUpdates(startSse, handler)

	return (
		<>
			<button className={"reset-button"} onClick={() => reset()}>Reset data</button>
			<div className="card">
				<div className="card-body">
					<input className="text-field" value={value} onChange={handleOnChange}/>
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
				</div>
				<div style={{display: "flex", justifyContent: "space-around"}}>
					<ul>
						<h3>Созданный данне</h3>
						{entity.map((el) => {
							return <Badge key={el.id} entity={el}/>
						})}
					</ul>
					<ul>
						<h3>Полученные данне</h3>
						{data?.map((el) => {
							return <Badge key={el.id} entity={el}/>
						})}
					</ul>
				</div>

			</div>
			<p className="read-the-docs">
				<button onClick={() => {
					create(entity).unwrap().then(() => {
						setTimeout(() => {
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
