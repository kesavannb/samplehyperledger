	/////////////////////////////////////////
	///////////// Setup Node.js /////////////
	/////////////////////////////////////////

		var express        = require('express');
		var app            = express();                // create our app w/ express
		var mongoose       = require('mongoose');      // mongoose for mongodb
		var morgan         = require('morgan');        // log requests to the console (express4)
		var bodyParser     = require('body-parser');     // pull information from HTML POST (express4)
		var methodOverride = require('method-override'); // simulate DELETE and PUT (express4)
		
		var inf;
	//// Set Server Parameters ////
	//var host = setup.SERVER.HOST;
	//var port = setup.SERVER.PORT;


	///////////  Configure app  ///////////

		//mongoose.connect('mongodb://node:nodeuser@mongo.onmodulus.net:27017/uwO3mypu');     // connect to mongoDB database on modulus.io

		app.use(express.static(__dirname + '/public'));                 // set the static files location /public/img will be /img for users
		app.use(morgan('dev'));                                         // log every request to the console
		app.use(bodyParser.urlencoded({'extended':'true'}));            // parse application/x-www-form-urlencoded
		app.use(bodyParser.json());                                     // parse application/json
		app.use(bodyParser.json({ type: 'application/vnd.api+json' })); // parse application/vnd.api+json as json
		app.use(methodOverride());
		/*var Todo = mongoose.model('Todo', {
			text : String
		});
		*/
		

	// ============================================================================================================================
	//				Handling request and response using post, get & delete method's
	// ============================================================================================================================

	//========================================
	//     Server details
	//========================================
	// listen (start app with node server.js) 
	//========================================
		app.listen(8080);
		console.log("App listening on port 8080");
		
	// ============================================================================================================================
	// 														Work Area
	// ============================================================================================================================
	var Ibc1 = require('ibm-blockchain-js');														//rest based SDK for ibm blockchain
	var ibc = new Ibc1();

	// ==================================================
	// configure options for ibm-blockchain-js sdk
	// ==================================================

	var chaincode = {}; 
	var ccdeployed = null;
	var parseval = null;

	var options = 	{
						network:{
							peers: [
									  {
										"api_host": "localhost", //replace with your hostname or ip of a peer                                        //replace with your https port (optional, omit if n/a)
										"api_port": 7050,        //replace with your http port
										"type": "peer",
										"id": "jdoe"             //unique id of peer
									  }
									],																	//lets only use the first peer! since we really don't need any more than 1
							users: [
									  {
										"enrollId": "bob",
										"enrollSecret": "NOE63pEQbL25"                        //enroll's secret
									  }
									],
						//dump the whole thing, sdk will parse for a good one
							options: {
										quiet: true, 															//detailed debug messages on/off true/false
										tls: false,
										//tls is used for enable the security purpose
										
	//									//should app to peer communication use tls?
										maxRetry: 1																//how many times should we retry register before giving up
									},								
						},
						chaincode:{
							
							git_url:'https://github.com/kesavannb/samplehyperledger/tree/master/samplechaincode',
							zip_url: 'https://github.com/kesavannb/samplehyperledger/archive/master.zip',
							//git_url:'https://github.com/bminchal/chain1/blob/master/git_code/',	
							//zip_url: 'https://github.com/bminchal/chain1/archive/master.zip',		
							//unzip_dir: '/chain1-master/git_code/',					
							unzip_dir:  'samplehyperledger-master/samplechaincode',
							 deployed_name: 'samplechaincode',
							
							//hashed cc name from prev deployment, comment me out to always deploy, uncomment me when its already deployed to skip deploying again
							//deployed_name: '16e655c0fce6a9882896d3d6d11f7dcd4f45027fd4764004440ff1e61340910a9d67685c4bb723272a497f3cf428e6cf6b009618612220e1471e03b6c0aa76cb'
						}
					};
							


	// ---- Fire off SDK ---- //



	//// Post method /////

	// create todo and send back all todos after creation
																																																					//sdk will populate this var in time, lets give it high scope by creating it here
	ibc.load(options, function (err, cc ){                                                              //parse/load chaincode, response has chaincode functions!
					if(err != null){
									console.log("options===>",options);
									
									console.log('! looks like an error loading the chaincode or network, app will fail\n', err);
					}
					else{
									chaincode = cc;                                 
									
									
	//Deploy is not working so i hide the line
									// ---- To Deploy or Not to Deploy ---- //
									if(!cc.details.deployed_name || cc.details.deployed_name === ''){                                                                         //yes, go deploy
													cc.deploy('init', ['99'], {delay_ms: 30000}, function(e){                                                                                    //delay_ms is milliseconds to wait after deploy for conatiner to start, 50sec recommended
																	check_if_deployed(e, 1);
													});
									}
									else{                                                                                                                                                                                                                                                                                                                      //no, already deployed
													console.log('chaincode summary file indicates chaincode has been previously deployed');
													check_if_deployed(null, 1);
									}	
					}
	});


	//loop here, check if chaincode is up and running or not

	function check_if_deployed(e, attempt){
					if(e){
									cb_deployed(e);                                                                                                                                                                                                                                                                                              //looks like an error pass it along
					}
					
					cb_deployed(null);
					
					
	}


	function cb_deployed(e){
					if(e != null){
									//look at tutorial_part1.md in the trouble shooting section for help
									console.log('! looks like a deploy error, holding off on the starting the socket\n', e);        
					}              
					else{
									console.log('------------------------------------------ Service Up ------------------------------------------');

									ccdeployed = "deployed";
									
					}
					
					
	}
		  



	//* Post code start 

	app.post('/api/insert', function(postreqang, postresang) {
													
	  // create a todo, information comes from AJAX request from Angular
					
		 var inf = postreqang.body;
		 console.log("TEXT FROM UI",inf);

			var uidata = inf;
									//var data = uidata.length;
					 //console.log("chaincode in data===>",data);
					 	 console.log("chaincode in uidata===>",uidata);
									var ID = uidata.ID;
									var Name  = uidata.Name; 
									var Details  = uidata.Details;
									var Index = uidata.Index;
					   
					   console.log("chaincode in invoke===>",chaincode);
					 
					  chaincode.invoke.save_data([ID,Name,Details,Index]);
					  
									  var response = postresang.end("inserted succesfully");
									  return response;
					  
					  
					
	  });

	app.post('/api/update', function(postreqang, postresang) {
													
	  // create a todo, information comes from AJAX request from Angular
					
		 var inf = postreqang.body;
		 console.log("TEXT FROM UI",inf);

			var uidata = inf;
									//var data = uidata;
					
									 var Index = uidata.Index;
									var Name  = uidata.Name; 
									var Details  = uidata.Details;
							  console.log("TEXT FROM UI data---",uidata);
					   console.log("chaincode in invoke===>",chaincode);
					
					  chaincode.invoke.update([Index, Name, Details]);
					  var response = postresang.end("Updated succesfully");
									  return response;
					
	  });
	  
	  app.delete('/api/delete', function(postreqang, postresang) {
													
	  // create a todo, information comes from AJAX request from Angular
					
		 var inf = postreqang.body;
		 console.log("TEXT FROM UI",inf);

			var uidata = inf;
									//var data = uidata;
					
									 var Index = uidata.Index;
									
							  console.log("TEXT FROM UI data---",uidata);
					   console.log("chaincode in invoke===>",chaincode);
					
					  chaincode.invoke.delete([Index]);
					  var response = postresang.end("Deleted succesfully");
									  return response;
					
	  });
	  
	  //* Post code end






		  
	//* Get code start 
									
	//// Get method /////

	app.get('/api/getAll', function(getreqang, getresang) {
			// use mongoose to get all todos in the database      
	chaincode.query.queryall([''], function(err, resp){                            
					
					if(resp != null){                  
									console.log("resp ===>",resp);
									var sTemp = "";
									var aObjs = [];
									
									for(var i=0; i<resp.length; ++i)
									{
													sTemp += resp[i];
													if (resp[i] == "}")
													{
																	aObjs.push(JSON.parse(sTemp));
																	sTemp = "";
																	//console.log("aObjs inside",aObjs);
													}
									} 

					console.log("aObjs", aObjs);                                       
					getresang.json(aObjs);
					//getresang.send(parseval2);
									  
					}                                              
	  });

 var response = getresang.end("Query all records succesfully");
									  return response;	  
	});

	app.get('/api/getByIndex/:id', function(getreqang, getresang) {
			// use mongoose to get all todos in the database      
										
		idValue = getreqang.params.id
									console.log("idValue",idValue);
	chaincode.query.query([idValue], function(err, resp){                   
					
					if(resp != null){                  
					console.log("resp ===>",resp); 
					var parseval = JSON.parse(resp);              
					console.log("parseval ==>",parseval);                                                    
					var info = {
									   "DATA": parseval
													   };                                          
					getresang.json(info);
									  
									}                                              
					}); 
 var response = getresang.end("Query based on the index id record succesfully");
									  return response;					
	});
					
					
	//* Get code end

	///// Get method //////               
	app.get('*', function(reqang1, resang1) {
		   resang1.sendfile('./public/index.html'); // load the single view file (angular will handle the page changes on the front-end)
	});

	/////// Executing post method from ui here ////////////



	  
	  //* Post code end

	// ========================================================
	// Block chain statics
	// ========================================================
			ibc.monitor_blockheight(function(chain_stats){										//there is a new block, lets refresh everything that has a state
				if(chain_stats && chain_stats.height){
					console.log('hey new block, lets refresh and broadcast to all', chain_stats.height-1);
					ibc.block_stats(chain_stats.height - 1);
					//chaincode.query.read(['_marbleindex'], cb_got_index);
				}
				
				});
			
	