/* Fallback for loading page */
const Fallback = () => {
	return (
		<div className="flex absolute inset-0 w-screen h-screen items-center justify-center bg-white z-[10001]">
			<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[15vw] min-w-[15rem] object-cover my-7 animate-pulse" draggable={false} />
		</div>
	);
};

export default Fallback;
