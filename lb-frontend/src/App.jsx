import React, { Fragment, useEffect, useState } from 'react'
import "axios";
import axios from 'axios';
const App = () => {
  const [subnet,setSubnet] = useState('');
  const [clientIp,setClientIp] = useState('');
  const [lbIp,setLbIp] = useState('');
  const [serverIp,setServerIp] = useState('');
  const handleSubnet = async (e) =>{
    e.preventDefault();
    const data = {
      subnet,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/create_subnet";
    const res = await axios.post(url,data,{headers});
    console.log(res.data);
  };
  const handleClientIp = async (e)=>{
    e.preventDefault();
    const data = {
      ip:clientIp,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/launch_client";
    const res = await axios.post(url,data,{headers});
    console.log(res.data);
  };
  const handleLbIp = async (e) => {
    e.preventDefault();
    const data = {
      ip:lbIp,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/launch_lb";
    const res = await axios.post(url,data,{headers});
    console.log(res.data);
  };
  const handleServerIp = async (e) => {
    e.preventDefault();
    const data = {
      ip:serverIp,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/launch_server";
    const res = await axios.post(url,data,{headers});
    console.log(res.data);
  };

  useEffect(()=>{
    const fetchInfo = async ()=> {
      console.log("fetching_data...");
    };
    fetchInfo();
  },[])
  return (
    <Fragment>
      <div className='container'>
          <div className='left-panel'>
            <div className='subnet-input'>
              <label>SUBNET </label>
              <input type="text" value={subnet} onChange={(e)=>{setSubnet(e.target.value)}}></input>
              <button onClick={(e)=> handleSubnet(e)}>Create Subnet</button>
            </div>
            <div className='clientip-input'>
              <label>CLIENT IP </label>
              <input type="text" value={clientIp} onChange={(e)=>{setClientIp(e.target.value)}}></input>
              <button onClick={(e)=> handleClientIp(e)}>Launch Client</button>
            </div>
            <div className='lb-input'>
              <label>LOAD-BALANCER IP </label>
              <input type="text" value={lbIp} onChange={(e)=>{setLbIp(e.target.value)}}></input>
              <button onClick={(e)=> handleLbIp(e)}>Launch LoadBalancer</button>
            </div>
            <div className='server-input'>
              <label>SERVER IP </label>
              <input type="text" value={serverIp} onChange={(e)=>{setServerIp(e.target.value)}}></input>
              <button onClick={(e)=> handleServerIp(e)}>Launch Server</button>
            </div>
          </div>
          <div className='right-panel'>
            {
            //TODO Table  key, ip-address, mac-address
            }
          </div>
      </div>
    </Fragment>
  )
}

export default App