// This is your test publishable API key.


document.addEventListener('DOMContentLoaded',async()=>{

    //fetch the publish stripe key
    const stripe = stripe("sk_test_51MrFKxSGSCYafJiXRlXzhLk1spPqCpRB6RnfUdykfhPu7MXzjg9QGyVhrrHlWQFYqOwzxKRK77uJGRqbCQYgySTN00CflxLkRk")

    const{clientSecret}=await fetch("/create-payment-intent",{
        method:"POST",
        headers:{
            "Content-Type": "application/json"
        },
    }).then(r=>r.json())

    //mount the elements
    const elements=stripe.elements({clientSecret})
    const paymentElement=elements.create('payment')
    paymentElement.mount('#payment-element')


})