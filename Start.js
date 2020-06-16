import React, { useState } from 'react';
import Button from '../button/Button';
import { Link } from 'react-router-dom';
import { getRoutes } from '../../clients/clients'
import './start.scss';


const Start = (props) => {
    const [busNumber, setBusNumber] = useState("")
    const [busData, setData] = useState([]);
    const { checkupContext } = props;
    const onClickConfirm = () => {
        checkupContext.setCheckupPayload({ ...checkupContext.checkupPayload, busNumber })
    }
    const queryBusData = () => {
        getRoutes(sessionStorage.getItem("accessToken"), busNumber)
            .then(res => {
                console.log("test:" + res)
                setData(res)
            })
            .catch(err => console.error(err))
    }
    return (
        <div className="content">
            <div className="start">
                <form onSubmit={(e) => {
                    e.preventDefault();
                    queryBusData();
                }}>
                    <input value={busNumber} className="start__input" type="text" placeholder="Įveskite transporto priemonės numerį" onChange={(e) => setBusNumber(e.target.value)} />
                    <button type="submit">send</button>
                </form>
                <div className="start__output">Pilnas transporto numeris autobusų parke</div>
                {
                    busData && busData.length > 0 ?
                        busData.map(data => <div>Bus Number: {
                            data.BusNumber
                        } Route: {data.Route}
                        </div>
                        )
                        :
                        null
                }
            </div>
            <Link to="/checkup/2" className="start-confirm">
                <Button className="btn btn-filled" onClick={() => onClickConfirm()} >Patvirtinti</Button>
            </Link>
        </div>
    )

}

export default Start;